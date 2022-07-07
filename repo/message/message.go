package message

import "gorm.io/gorm"

type MessageRepo struct {
	DB *gorm.DB
}

func NewMessageRepo(db *gorm.DB) MessageRepoInterface {
	return &MessageRepo{
		DB: db,
	}
}

func (r *MessageRepo) Save(m *MessageModel) {
	r.DB.Create(&m)
}
