package redis

import (
	"time"

	redis "github.com/go-redis/redis/v8"
)

type IRedis interface {
	Ping() (string, error)
	Set(string, interface{}, time.Duration) error
	Get(string) (interface{}, error)
	GetAllKeys() ([]string, error)
	Delete(string) error
	Nil() error
	Subscribe(string) (*redis.PubSub, error)
	Publish(string, interface{}) error
	Close() error
}
