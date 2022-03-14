package redis

import (
	"context"
	"fmt"
	"time"

	redis "github.com/go-redis/redis/v8"
)

type Redis struct {
	client *redis.Client
	ctx    context.Context
	pubsub redis.PubSub
}

func New(host string, port int, password string, DB int, ctx context.Context) (IRedis, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       DB,
	})

	fmt.Println(fmt.Sprintf("%s:%d", host, port))
	_, err := rdb.Ping(ctx).Result()
	return &Redis{
		client: rdb,
		ctx:    ctx,
	}, err
}

func (c *Redis) Ping() (string, error) {
	return c.client.Ping(c.ctx).Result()
}

func (c *Redis) Set(key string, value interface{}, expiration time.Duration) error {
	err := c.client.Set(c.ctx, key, value, expiration).Err()
	return err
}

func (c *Redis) Get(key string) (interface{}, error) {
	value, err := c.client.Get(c.ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf(fmt.Sprintf("Key `%s` does not exist", key))
	}
	return value, err
}

func (c *Redis) Nil() error {
	return redis.Nil
}

func (c *Redis) GetAllKeys() ([]string, error) {
	resultArr := []string{}
	keys, err := c.client.Do(c.ctx, "KEYS", "*").Result()

	// we have to do manual conversion
	// due to Golang strict typing
	for _, val := range keys.([]interface{}) {
		resultArr = append(resultArr, val.(string))
	}

	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("An error occured: %v", err))
	}
	return resultArr, err
}

func (c *Redis) Delete(key string) error {
	err := c.client.Del(c.ctx, key).Err()
	if err == redis.Nil {
		return fmt.Errorf(fmt.Sprintf("Key `%s` does not exist", key))
	}
	return err
}

func (c *Redis) Subscribe(channel string) (*redis.PubSub, error) {
	pubsub := c.client.Subscribe(c.ctx, channel)
	iface, err := pubsub.Receive(c.ctx)
	switch iface.(type) {
	case *redis.Subscription:
	case *redis.Message:
		fmt.Println("Subscription is succeeded")
	}
	return pubsub, err
}

func (c *Redis) Publish(channel string, value interface{}) error {
	err := c.client.Publish(c.ctx, channel, value).Err()
	return err
}

func (c *Redis) Close() error {
	err := c.client.Close()
	return err
}
