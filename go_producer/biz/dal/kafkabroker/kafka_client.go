package kakfabroker

import (
	"context"
	"encoding/json"
	"fmt"
	"lintang/go_producer/biz/domain"
	"lintang/go_producer/config"
	"math/rand"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type KafkaProducer struct {
	Writer *kafka.Writer
}

func NewKafkaProducer(cfg *config.Config) *KafkaProducer {
	dialer := &kafka.Dialer{
		Timeout:  10 * time.Second,
		ClientID: "oti-be" + generateRandomString(5),
	}

	config := kafka.WriterConfig{
		Brokers:      []string{cfg.Kafka.KafkaAddress},
		Topic:        cfg.Kafka.NotificationTopic,
		Balancer:     &kafka.Murmur2Balancer{},
		Dialer:       dialer,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	w := kafka.NewWriter(config)
	zap.L().Info(fmt.Sprintf("successfully connected to kafka broker %s as producer!!", cfg.Kafka.KafkaAddress))

	return &KafkaProducer{w}
}

func (p KafkaProducer) Publish(ctx context.Context, key, value []byte) error {
	message := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}

	return p.Writer.WriteMessages(ctx, message)
}

func (p KafkaProducer) PublishNotificationsToKafka(ctx context.Context, key string, notifMsg domain.NotificationMessage) error {

	var msgInBytes []byte
	msgInBytes, err := json.Marshal(notifMsg)
	if err != nil {
		return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}

	keyInBytes, err := json.Marshal(key)
	if err != nil {
		return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}

	if err := p.Publish(ctx, keyInBytes, msgInBytes); err != nil {
		return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)

	}
	zap.L().Info(fmt.Sprintf("publishing message %s to kafka!!", notifMsg.Message))

	return nil
}

func (p KafkaProducer) Close(ctx context.Context) {
	p.Writer.Close()
}



func generateRandomString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
