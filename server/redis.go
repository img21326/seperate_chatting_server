package server

import "github.com/go-redis/redis/v8"

func InitRedis(option *redis.Options) *redis.Client {
	if option == nil {
		option = &redis.Options{
			Addr:     "139.162.125.28:6379",
			Password: "",
			DB:       5,
		}
	}
	client := redis.NewClient(option)
	return client
}
