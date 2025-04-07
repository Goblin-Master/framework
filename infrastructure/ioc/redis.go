package ioc

import (
	"context"
	"ddd/infrastructure/config"
	"fmt"
	"github.com/go-redis/redis/v8"
)

func InitRedis() redis.Cmdable {
	if !config.Conf.Redis.Enable {
		return nil
	}
	client := redis.NewClient(&redis.Options{
		Network:            "",
		Addr:               config.Conf.Redis.DSN(),
		Dialer:             nil,
		OnConnect:          nil,
		Username:           "",
		Password:           config.Conf.Redis.Password,
		DB:                 config.Conf.Redis.DB,
		MaxRetries:         0,
		MinRetryBackoff:    0,
		MaxRetryBackoff:    0,
		DialTimeout:        0,
		ReadTimeout:        0,
		WriteTimeout:       0,
		PoolFIFO:           false,
		PoolSize:           1000,
		MinIdleConns:       1,
		MaxConnAge:         0,
		PoolTimeout:        0,
		IdleTimeout:        0,
		IdleCheckFrequency: 0,
		TLSConfig:          nil,
		Limiter:            nil,
	})
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		panic(fmt.Sprintf("redis init error:%v", err))
	}
	return client
}
