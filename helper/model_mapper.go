package helper

import (
	"kama_chat_server/dto"
	"kama_chat_server/model"
	"time"

	"github.com/google/uuid"
)

func ChatRequest2Message(source *dto.ChatRequest) *model.Message {
	return &model.Message{
		Uuid:       uuid.New().String(),
		SessionId:  source.SessionId,
		Type:       source.Type,
		Content:    source.Content,
		Url:        "",
		SendId:     source.SendId,
		SendName:   source.SendName,
		SendAvatar: source.SendAvatar,
		ReceiveId:  source.ReceiveId,
		FileSize:   "0B",
		FileType:   "",
		FileName:   "",
		Status:     model.UNSENT,
		CreatedAt:  time.Now(),
		AVdata:     "",
	}
}

func Message2ChatResponse(source *model.Message) *dto.ChatResponse {
	return &dto.ChatResponse{
		Uuid:       source.Uuid,
		SendId:     source.SendId,
		SendName:   source.SendName,
		SendAvatar: source.SendAvatar,
		ReceiveId:  source.ReceiveId,
		Type:       source.Type,
		Content:    source.Content,
		Url:        source.Url,
		FileSize:   source.FileSize,
		FileName:   source.FileName,
		FileType:   source.FileType,
		CreatedAt:  source.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
