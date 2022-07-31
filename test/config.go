package test

import (
	"github.com/glebarez/sqlite"
	"github.com/go-redis/redis/v8"
	"github.com/img21326/fb_chat/server"
)

var GormDialector = sqlite.Open(":memory:")

var RedisOption = &redis.Options{
	Addr:     "139.162.125.28:6379",
	Password: "",
	DB:       5,
}

var DB = server.InitDB(GormDialector)
var Redis = server.InitRedis(RedisOption)
var URL = "http://localhost"
