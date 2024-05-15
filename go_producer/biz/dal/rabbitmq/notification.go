package rabbitmq

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"lintang/go_producer/biz/domain"
	"time"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type NotificationMQ struct {
	ch *amqp.Channel
}

func NewNotificationMQ(rmq *RabbitMQ) *NotificationMQ {
	return &NotificationMQ{
		rmq.Channel,
	}
}

func (n *NotificationMQ) SendEmailNotificationsToRMQ(ctx context.Context, msg domain.NotificationMessage) error {
	zap.L().Info(fmt.Sprintf("sending message %s to rabbit mq!! ", string(msg.Message)))

	return n.publish(ctx, "go-notifications", "go-email-notification", msg)
}

func (m *NotificationMQ) publish(ctx context.Context, exchange string, routingKey string, event interface{}) error {
	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(event); err != nil {
		zap.L().Error("gob.NewEncoder(&b).Encode(event)", zap.Error(err))
		return err
	}

	err := m.ch.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,
		false,
		amqp.Publishing{
			AppId:       "producer-rabbitmq-go-server",
			ContentType: "application/x-encoding-gob",
			Body:        b.Bytes(),
			Timestamp:   time.Now(),
		})
	if err != nil {
		zap.L().Error("m.ch.Publish: ", zap.Error(err))
		return err
	}

	return nil
}
