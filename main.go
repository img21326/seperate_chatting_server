package main

import (
	"log"
	"os"

	"github.com/img21326/fb_chat/helper"
	"github.com/img21326/fb_chat/server"
	"github.com/spf13/viper"
)

func main() {

	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	db := server.InitDB(nil)
	redis := server.InitRedis(nil)
	port := helper.GetEnv("PORT", "8080")

	server.StartUpRedisServer(db, redis, port)
}
