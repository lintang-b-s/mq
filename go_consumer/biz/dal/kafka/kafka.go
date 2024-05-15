package kafka

import (
	"context"
	"fmt"
	"lintang/go_consumer/config"
	"math/rand"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type KafkaConsumerGroup struct {
	Reader *kafka.Reader
}

func NewKafkaConsumerGroup(cfg *config.Config) *KafkaConsumerGroup {
	groupID := "go-consumer-group" + generateRandomString(4)
	config := kafka.ReaderConfig{
		Brokers: []string{cfg.Kafka.KafkaAddress},
		GroupID: groupID,
		Topic:   cfg.Kafka.NotificationTopic,
	}
	reader := kafka.NewReader(config)
	zap.L().Info(fmt.Sprintf("successfully connected to kafka broker %s  as consumer group %s", cfg.Kafka.KafkaAddress, groupID))
	return &KafkaConsumerGroup{reader}
}

func (c KafkaConsumerGroup) ReadNotificationMessage() {

	go func() {
		for {
			m, err := c.Reader.ReadMessage(context.Background())
			if err != nil {
				zap.L().Error("c.Reader.ReadMessage (ReadNotificationMessage)", zap.Error(err))
				continue
			}
			topic := m.Topic
			value := m.Value

			if err != nil {
				zap.L().Error("errror while receiving mesages", zap.Error(err))
				continue
			}

			zap.L().Info(fmt.Sprintf("consumed message %s from topic %v", string(value), topic))
		}
	}()
}

func (c KafkaConsumerGroup) Close(ctx context.Context) {
	c.Reader.Close()
}

func generateRandomString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
