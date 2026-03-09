package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"kama_chat_server/kafka"
)

func SendMessage(id string, mess string) {

}

func ReceiveMessage(ctx context.Context) (*KafkaRequest, error) {
	data, err := kafka.ReadMessage(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read message from kafka: %v", err)
	}
	var req KafkaRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, fmt.Errorf("failed to parse request from kafka: %v", err)
	}
	return &req, nil
}
