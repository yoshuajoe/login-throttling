package main

import (
	"context"
	"fmt"
	"login-throttling/internal/pkg/redis"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	BUCKET    = 1 * 60
	THRESHOLD = 10
)

type LogMiddleware struct {
}

func NewLogMiddleware() LogMiddleware {
	return LogMiddleware{}
}

func GetKey(IP string) string {
	bucket := time.Now().Unix() / BUCKET
	IP = IP + strconv.FormatInt(bucket, 10)
	return IP
}

func RGet(key string, r redis.IRedis) ([]string, error) {
	rGet, rGetErr := r.Get(key)

	if rGet != nil {
		if rGetErr != r.Nil() && rGet != nil {
			return strings.Split(rGet.(string), ","), rGetErr
		}

		return []string{
			rGet.(string),
		}, rGetErr
	}
	return []string{}, nil
}

func (l LogMiddleware) LogWriter(authSecret string, host string, port int, auth string, duration int, ctx context.Context, cancel context.CancelFunc) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			redisClient, redisClientErr := redis.New(host, port, auth,
				0, ctx)
			defer redisClient.Close()

			if redisClientErr != nil {
				cancel()
				panic(redisClientErr)
			}

			IPAddress := c.Request().Header.Get("X-Real-Ip")
			if IPAddress == "" {
				IPAddress = c.Request().Header.Get("X-Forwarded-For")
			}
			if IPAddress == "" {
				IPAddress = c.Request().RemoteAddr
			}

			modifiedIPAddress := GetKey(IPAddress)
			fmt.Println("IP:", modifiedIPAddress)

			val, _ := RGet(modifiedIPAddress, redisClient)

			if val == nil || len(val) <= 0 {
				redisClient.Set(modifiedIPAddress, 1, time.Second*time.Duration(duration))
				fmt.Println("Created")
			} else {
				valInt, _ := strconv.ParseInt(val[0], 10, 0)

				if valInt+1 > THRESHOLD {
					return c.String(401,
						"Max retries exceeded, please wait "+strconv.Itoa(duration)+" seconds",
					)
				}
				redisClient.Set(modifiedIPAddress, valInt+1, time.Second*time.Duration(duration))
			}

			authRaw := c.Request().Header.Get("X-SERVER-SECRET-TOKEN")
			if authRaw == authSecret {
				return next(c)
			}

			return c.NoContent(401)
		}
	}
}
