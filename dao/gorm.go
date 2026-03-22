package dao

import (
	"fmt"
	"kama_chat_server/config"
	"kama_chat_server/model"
	"kama_chat_server/zlog"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	cfg := config.GetConfig().MysqlConfig
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	var err error
	db, err = gorm.Open(mysql.Open(dsn))
	if err != nil {
		zlog.Fatal("failed to open mysql: %v", err)
	}
	if err := db.AutoMigrate(&model.Message{}); err != nil {
		zlog.Fatal("failed to auto migrate: %v", err)
	}
}

func Insert(m any) error {
	tx := db.Create(m)
	return tx.Error
}

func GetMessageList(userOneId, userTwoId string) ([]*model.Message, error) {
	var messages []*model.Message
	res := db.
		Where("(send_id = ? AND receive_id = ?) OR (send_id = ? AND receive_id = ?)", userOneId, userTwoId, userTwoId, userOneId).
		Order("created_at ASC").
		Find(&messages)
	if res.Error != nil {
		return nil, fmt.Errorf("failed to select message list from db: %v", res.Error)
	}
	return messages, nil
}

func GetGroupMessageList(groupId string) ([]*model.Message, error) {
	var messages []*model.Message
	res := db.
		Where("receive_id = ?", groupId).
		Order("created_at ASC").
		Find(&messages)
	if res.Error != nil {
		return nil, fmt.Errorf("failed to select group message list from db: %v", res.Error)
	}
	return messages, nil
}
