package main

import (
	"fmt"
	"log"
	"os"

	"github.com/img21326/fb_chat/server"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	var db *gorm.DB
	if viper.GetString("POSTGRES_HOST") != "" {
		dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Taipei",
			viper.GetString("POSTGRES_HOST"), viper.GetString("POSTGRES_USER"), viper.GetString("POSTGRES_PWD"), viper.GetString("POSTGRES_DB"), viper.GetString("POSTGRES_PORT"))
		dbConfig := postgres.Open(dsn)
		db = server.InitDB(dbConfig)
	} else {
		db = server.InitDB(nil)
	}

	port := viper.GetString("SERVER_PORT")

	if viper.GetString("SERVER_TYPE") == "REDIS_SERVER" {
		redis := server.InitRedis(nil)
		server.StartUpRedisServer(db, redis, port)
	} else if viper.GetString("SERVER_TYPE") == "LOCAL_SERVER" {
		server.StartUpLocalServer(db, port)
	} else {
		fmt.Print("ENV SET ERROR")
	}

}
