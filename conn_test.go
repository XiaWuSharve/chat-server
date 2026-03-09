package mychat

import (
	"context"
	"kama_chat_server/dao"
	"kama_chat_server/dto"
	"kama_chat_server/helper"
	"kama_chat_server/kafka"
	"testing"
)

func TestKafka(t *testing.T) {
	kafka.Start()
	ctx := context.Background()
	err := kafka.SendMessage(ctx, "hello kafka")
	if err != nil {
		t.Error(err)
	}
	data, err := kafka.ReadMessage(ctx)
	if err != nil {
		t.Error(err)
	}
	if string(data) != "hello kafka" {
		t.Errorf("expected %s, got %s", "hello kafka", string(data))
	}
}

func TestMysql(t *testing.T) {
	m := helper.KafkaRequest2Message(&dto.KafkaRequest{
		Content: "hello mysql",
	})
	if err := dao.Insert(m); err != nil {
		t.Error(err)
	}
}
