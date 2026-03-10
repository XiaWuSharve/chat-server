package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"kama_chat_server/config"
	"kama_chat_server/dto"
	"kama_chat_server/zlog"
	"time"

	"github.com/go-redis/redis/v8"
)

var client *redis.Client
var expire time.Duration

func init() {
	cfg := config.GetConfig().RedisConfig
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.Db,
	})
	expire = cfg.Expire
}

func get(ctx context.Context, key string) (string, error) {
	val, err := client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			zlog.Info("cache miss: %s", key)
		}
		return "", err
	}
	return val, nil
}

func SetPrivateMessage(ctx context.Context, sendId string, receiveId string, messages []*dto.ChatResponse) error {
	data, err := json.Marshal(messages)
	if err != nil {
		return fmt.Errorf("failed to parse message list from %s to %s: %v", sendId, receiveId, err)
	}
	var key string
	if sendId < receiveId {
		key = fmt.Sprintf("messages:%s:%s", sendId, receiveId)
	} else {
		key = fmt.Sprintf("messages:%s:%s", receiveId, sendId)
	}
	return set(ctx, key, string(data), expire*time.Minute)
}

func GetPrivateMessage(ctx context.Context, sendId string, receiveId string) ([]*dto.ChatResponse, error) {
	var key string
	if sendId < receiveId {
		key = fmt.Sprintf("messages:%s:%s", sendId, receiveId)
	} else {
		key = fmt.Sprintf("messages:%s:%s", receiveId, sendId)
	}
	s, err := get(ctx, key)
	if err != nil {
		if errors.Is(redis.Nil, err) {
			s = "[]"
		} else {
			return nil, fmt.Errorf("failed to get key %s: %v", key, err)
		}
	}
	var req []*dto.ChatResponse
	if err := json.Unmarshal([]byte(s), &req); err != nil {
		return nil, fmt.Errorf("failed to parse redis value %s: %v", s, err)
	}
	return req, nil
}

func set(ctx context.Context, key string, value string, expire time.Duration) error {
	err := client.Set(ctx, key, value, expire).Err()
	return err
}
