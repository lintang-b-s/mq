package rabbitmq

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"lintang/go_consumer/biz/domain"

	"go.uber.org/zap"
)

type NotificationListener struct {
	rmq *RabbitMQ

	done chan struct{}
}

func NewNotificationListener(rmq *RabbitMQ, d chan struct{}) *NotificationListener {
	return &NotificationListener{rmq, d}
}

func (l NotificationListener) ListenAndServe() error {
	queue, err := l.rmq.Channel.QueueDeclare(
		"go-email-notification-queue",
		false,
		false,
		true,
		false,
		nil,
	)

	if err != nil {
		return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}

	err = l.rmq.Channel.QueueBind(
		queue.Name,
		"go-email-notification",
		"go-notifications",
		false,
		nil,
	)

	if err != nil {
		return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}

	msgs, err := l.rmq.Channel.Consume(
		queue.Name,
		"go-email-consumer",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}

	var nack bool
	go func() {
		for msg := range msgs {
			zap.L().Info(fmt.Sprintf(`received message: %s`, msg.RoutingKey))

			switch msg.RoutingKey {
			case "go-email-notification":
				notification, err := decodeNotificationMessage(msg.Body)
				if err != nil {
					zap.L().Error("decodeMetadataMessage (ListenAndServe)", zap.Error(err))
				}

				zap.L().Info(fmt.Sprintf("consumed message '%s' from rabbitmq!!", string(notification.Message)))

				if err != nil {
					nack = true
				}
			}

			if nack {
				zap.L().Info("nack")
				_ = msg.Nack(false, nack)
			} else {
				zap.L().Info("ack")
				_ = msg.Ack(false)
			}
		}

		l.done <- struct{}{}
	}()

	return nil

}

type NotificationMesage struct {
	Message string `json:"message"`
}

func decodeNotificationMessage(b []byte) (NotificationMesage, error) {
	var res NotificationMesage
	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(&res); err != nil {
		zap.L().Error("NewDecoder (decodeMetadataMessage) (MetadataListener)", zap.Error(err))
		return NotificationMesage{}, domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}
	return res, nil
}
