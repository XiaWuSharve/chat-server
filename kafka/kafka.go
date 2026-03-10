package kafka

import (
	"context"
	"fmt"
	"kama_chat_server/config"
	"kama_chat_server/zlog"
	"time"

	"github.com/segmentio/kafka-go"
)

var (
	chatWriter *kafka.Writer
	chatReader *kafka.Reader
	cfg        = config.GetConfig().KafkaConfig
)

func Start() {
	createTopicIfNotExists(cfg.ChatTopic)
	chatWriter = &kafka.Writer{
		Addr:                   kafka.TCP(cfg.Address),
		Topic:                  cfg.ChatTopic,
		Balancer:               &kafka.Hash{},
		WriteTimeout:           cfg.Timeout * time.Second,
		RequiredAcks:           kafka.RequireNone,
		AllowAutoTopicCreation: false,
	}
	chatReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{cfg.Address},
		Topic:          cfg.ChatTopic,
		CommitInterval: cfg.Timeout * time.Second,
		GroupID:        "chat",
		StartOffset:    kafka.FirstOffset,
	})
}

func Close() error {
	if err := chatWriter.Close(); err != nil {
		return fmt.Errorf("failed to close chat writer: %v", err)
	}
	if err := chatReader.Close(); err != nil {
		return fmt.Errorf("failed to close chat reader: %v", err)
	}
	return nil
}

func SendMessage(ctx context.Context, data []byte) error {
	return chatWriter.WriteMessages(ctx, kafka.Message{
		Value: data,
	})
}

func ReadMessage(ctx context.Context) ([]byte, error) {
	m, err := chatReader.ReadMessage(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read message from kafka: %v", err)
	}
	zlog.Info(fmt.Sprintf("topic=%s, partition=%d, offset=%d, key=%s, value=%s", m.Topic, m.Partition, m.Offset, m.Key, m.Value))
	return m.Value, nil
}

func createTopicIfNotExists(topic string) {
	conn, err := kafka.Dial("tcp", cfg.Address)
	if err != nil {
		zlog.Fatal(fmt.Errorf("failed to connect to kafka: %v", err))
	}
	defer conn.Close()
	err = conn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     cfg.Partition,
		ReplicationFactor: 1,
	})
	if err != nil {
		zlog.Fatal(err)
	}
}
