package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"login-throttling/internal/app"
	"login-throttling/internal/pkg/config"

	"github.com/caarlos0/env"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func readConfig(filename string, conf *config.Config) error {
	appsConf := config.SApps{}
	redisConf := config.SRedis{}

	appsConfParsedErr := env.Parse(&appsConf)
	if appsConfParsedErr != nil {
		return appsConfParsedErr
	}

	redisConfParsedErr := env.Parse(&redisConf)
	if redisConfParsedErr != nil {
		return redisConfParsedErr
	}

	conf.Apps = appsConf
	conf.Redis = redisConf
	return nil
}

func main() {
	defer recovery(5, logic)
}

func logic() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	conf := config.Config{}

	// Read the config yaml file
	err := readConfig("config", &conf)
	if err != nil {
		cancel()
		panic(err)
	}

	app, appErr := app.New()
	if appErr != nil {
		cancel()
		panic(appErr)
	}

	logMiddleware := NewLogMiddleware()

	// routes
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		Skipper:           middleware.DefaultSkipper,
		StackSize:         4 << 10, // 4 KB
		DisableStackAll:   true,
		DisablePrintStack: true,
		LogLevel:          0,
	}))
	e.Use(logMiddleware.LogWriter(conf.Apps.BasicAuthStatic, conf.Redis.Host, conf.Redis.Port, conf.Redis.Auth, conf.Redis.CacheExpiry, ctx, cancel))

	e.POST("/login", app.Login)
	// end routes

	// handle SIGTERM
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go sigtermHandler(c, cancel)
	go e.Logger.Fatal(e.Start(":1328"))
}

func recovery(maxRetry int, f func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(
				fmt.Sprintf("Panic occured %v, recovery in progress", r),
			)
			if maxRetry == 0 {
				panic("Too many PANIC retried")
			} else {
				time.Sleep(5 * time.Second)
				recovery(maxRetry-1, f)
			}
		}
	}()
	f()
}

func sigtermHandler(c chan os.Signal, canc context.CancelFunc) {
	<-c
	defer os.Exit(1)
	canc()
}
