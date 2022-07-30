package server

import (
	"log"

	"github.com/glebarez/sqlite"
	ModelMessage "github.com/img21326/fb_chat/structure/message"
	ModelRoom "github.com/img21326/fb_chat/structure/room"
	ModelUser "github.com/img21326/fb_chat/structure/user"
	"gorm.io/gorm"
)

func InitDB(dialector gorm.Dialector) *gorm.DB {
	if dialector == nil {
		dialector = sqlite.Open("gorm.db")
	}
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Panicf("Open db error: %v", err)
	}
	db.AutoMigrate(&ModelUser.User{}, &ModelMessage.Message{}, &ModelRoom.Room{})
	return db
}
