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
		zlog.Fatal(err)
	}
	if err := db.AutoMigrate(&model.Message{}); err != nil {
		zlog.Fatal(err)
	}
}

func Insert(m any) error {
	tx := db.Create(m)
	return tx.Error
}
