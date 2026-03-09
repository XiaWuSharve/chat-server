package helper

import (
	"kama_chat_server/dto"
	"kama_chat_server/model"
	"time"

	"github.com/google/uuid"
)

func KafkaRequest2Message(source *dto.KafkaRequest) *model.Message {
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
