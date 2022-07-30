package main

import (
	"os"

	"github.com/img21326/fb_chat/server"
)

func main() {

	// if helper.GetEnv("HOST_NAME", "") == "" {
	// 	panic("Please set env HOST_NAME")
	// }

	db := server.InitDB(nil)
	redis := server.InitRedis(nil)
	port := os.Args[1]

	server.StartUpRedisServer(db, redis, port)
}
