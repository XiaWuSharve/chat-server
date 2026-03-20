package config

import (
	"fmt"
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

type RedisConfig struct {
	Host     string        `toml:"host"`
	Port     int           `toml:"port"`
	Password string        `toml:"password"`
	Db       int           `toml:"db"`
	Expire   time.Duration `toml:"expire"` // in minute
}

type LogConfig struct {
	LogPath string `toml:"log_path"`
}

type Config struct {
	KafkaConfig *KafkaConfig `toml:"kafka"`
	MysqlConfig *MysqlConfig `toml:"mysql"`
	RedisConfig *RedisConfig `toml:"redis"`
	LogConfig   *LogConfig   `toml:"log"`
}

var config = &Config{
	KafkaConfig: &KafkaConfig{
		Address:     "127.0.0.1:9092",
		LoginTopic:  "login",
		ChatTopic:   "chat_message",
		LogoutTopic: "logout",
		Timeout:     1,
		Partition:   1,
	},
	MysqlConfig: &MysqlConfig{
		User:     "root",
		Password: "root",
		Host:     "127.0.0.1",
		Port:     3306,
		Database: "sharve",
	},
	RedisConfig: &RedisConfig{
		Host:     "127.0.0.1",
		Port:     6379,
		Password: "",
		Db:       0,
		Expire:   1440,
	},
	LogConfig: &LogConfig{
		LogPath: "./production.log",
	},
}

func LoadConfig() error {
	_, err := toml.DecodeFile("./configs/config.toml", &config)
	return err
}

func GetConfig() *Config {
	if config == nil {
		panic(fmt.Errorf("must LoadConfig first!"))
	}
	return config
}
