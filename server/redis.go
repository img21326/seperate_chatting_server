package server

import (
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

func InitRedis(option *redis.Options) *redis.Client {
	vip := viper.GetViper()

	if option == nil {
		option = &redis.Options{
			Addr:     vip.GetString("REDIS_HOST"),
			Password: vip.GetString("REDIS_PWD"),
			DB:       vip.GetInt("REDIS_DB"),
		}
	}
	client := redis.NewClient(option)
	return client
}
