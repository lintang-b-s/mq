package service

import (
	"context"
	"lintang/go_producer/biz/domain"
)

type KafkaProducer interface {
	PublishNotificationsToKafka(ctx context.Context, key string, notifMsg domain.NotificationMessage) error
}

type KafkaProducerService struct {
	kafka KafkaProducer
}

func NewKafkaProducerService(kafka KafkaProducer) *KafkaProducerService {
	return &KafkaProducerService{
		kafka,
	}
}

func (p *KafkaProducerService) SendEmailNotificationsToKafka(ctx context.Context, message string) error {
	err := p.kafka.PublishNotificationsToKafka(ctx, "go-kafka-notifications-key", domain.NotificationMessage{
		Message: message,
	})
	if err != nil {
		return err
	}
	return nil
}
