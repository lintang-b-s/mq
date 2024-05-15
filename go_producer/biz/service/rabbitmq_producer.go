package service

import (
	"context"
	"lintang/go_producer/biz/domain"
)

type NotificationMQ interface {
	SendEmailNotificationsToRMQ(ctx context.Context, msg domain.NotificationMessage) error
}

type RabbitMQProducerService struct {
	rmq NotificationMQ
}

func NewRabbitMQProducerService(rmq NotificationMQ) *RabbitMQProducerService {
	return &RabbitMQProducerService{
		rmq,
	}
}

func (p *RabbitMQProducerService) SendEmailNotificationsToRMQ(ctx context.Context, message string) error {
	err := p.rmq.SendEmailNotificationsToRMQ(ctx, domain.NotificationMessage{
		Message: message,
	})
	if err != nil {
		return err
	}
	return nil
}
