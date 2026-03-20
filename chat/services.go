package chat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"kama_chat_server/dao"
	"kama_chat_server/dto"
	"kama_chat_server/helper"
	"kama_chat_server/kafka"
	"kama_chat_server/redis"
	"kama_chat_server/zlog"
)

func SendMessage(id string, mess string) {

}

func KafkaReceiveMessage(ctx context.Context) (*dto.ChatRequest, error) {
	data, err := kafka.ReadMessage(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read message from kafka: %v", err)
	}
	var req dto.ChatRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, fmt.Errorf("failed to parse request from kafka: %v", err)
	}
	return &req, nil
}

func KafkaSendMessage(ctx context.Context, req *dto.ChatRequest) error {
	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}
	if err = kafka.SendMessage(ctx, data); err != nil {
		return fmt.Errorf("failed to send message to kafka: %v", err)
	}
	return nil
}

func RedisAddPrivateMessage(ctx context.Context, mess *dto.ChatResponse) error {
	messages, err := redis.GetPrivateMessage(ctx, mess.SendId, mess.ReceiveId)
	if err != nil {
		return fmt.Errorf("failed to get a message from redis: %v", err)
	}
	messages = append(messages, mess)

	if err = redis.SetPrivateMessage(ctx, mess.SendId, mess.ReceiveId, messages); err != nil {
		return fmt.Errorf("failed to set a message to redis: %v", err)
	}
	return nil
}

func GetMessageList(ctx context.Context, userOneId, userTwoId string) (message string, rsp []*dto.ChatResponse, ret int) {
	messages, err := redis.GetPrivateMessage(ctx, userOneId, userTwoId)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			zlog.Error("failed to hit redis for message about %s and %s", userOneId, userTwoId)
			msgs, err := dao.GetMessageList(userOneId, userTwoId)
			if err != nil {
				return "数据库访问失败", nil, -1
			}
			// TODO: set to redis
			messages = helper.BatchOperation(msgs, helper.Message2ChatResponse)
			return "成功获取消息列表", messages, 0
		}
		return "消息列表获取失败", nil, -1
	}
	return "成功获取消息列表", messages, 0
}
