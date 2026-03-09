package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"kama_chat_server/dto"
	"kama_chat_server/kafka"
	"kama_chat_server/redis"
)

func SendMessage(id string, mess string) {

}

func KafkaReceiveMessage(ctx context.Context) (*dto.KafkaRequest, error) {
	data, err := kafka.ReadMessage(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read message from kafka: %v", err)
	}
	var req dto.KafkaRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, fmt.Errorf("failed to parse request from kafka: %v", err)
	}
	return &req, nil
}

func RedisAddPrivateMessage(ctx context.Context, mess *dto.RedisMessage) error {
	messages, err := redis.GetPrivateMessage(ctx, mess.SendId, mess.ReceiveId)
	if err != nil {
		return fmt.Errorf("failed to get a message to redis: %v", err)
	}
	messages = append(messages, mess)

	if err = redis.SetPrivateMessage(ctx, mess.SendId, mess.ReceiveId, messages); err != nil {
		return fmt.Errorf("failed to set a message to redis: %v", err)
	}
	return nil
}
