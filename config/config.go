package config

import (
	"log"
	"time"

	"github.com/BurntSushi/toml"
)

type KafkaConfig struct {
	Address     string        `toml:"address"`
	LoginTopic  string        `toml:"login_topic"`
	ChatTopic   string        `toml:"chat_topic"`
	LogoutTopic string        `toml:"logout_topic"`
	Timeout     time.Duration `toml:"timeout"`
	Partition   int           `toml:"partition"`
}

type MysqlConfig struct {
	User     string `toml:"user"`
	Password string `toml:"password"`
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Database string `toml:"database"`
}

type Config struct {
	KafkaConfig *KafkaConfig `toml:"kafka"`
	MysqlConfig *MysqlConfig `toml:"mysql"`
}

var config = &Config{
	KafkaConfig: &KafkaConfig{
		Address:     "127.0.0.1:9092",
		LoginTopic:  "login",
		ChatTopic:   "chat_message",
		LogoutTopic: "logout",
		Timeout:     1,
	},
	MysqlConfig: &MysqlConfig{
		User:     "root",
		Password: "root",
		Host:     "127.0.0.1",
		Port:     3306,
		Database: "sharve",
	},
}

func init() {
	_, err := toml.DecodeFile("./configs/config.toml", &config)
	if err != nil {
		log.Fatal(err)
	}
}

func GetConfig() *Config {
	return config
}
