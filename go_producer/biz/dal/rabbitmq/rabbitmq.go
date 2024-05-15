package rabbitmq

import (
	"context"
	"fmt"
	"lintang/go_producer/config"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func NewRabbitMQ(cfg *config.Config) *RabbitMQ {
	conn, err := amqp.Dial(cfg.RabbitMQ.RMQAddress)
	if err != nil {
		hlog.Fatal("error: cannot connect to rabbitmq: " + err.Error())
	}

	channel, err := conn.Channel()
	if err != nil {
		hlog.Fatal("error can't get rabbitmq cahnnel: " + err.Error())
	}

	err = channel.ExchangeDeclare(
		"go-notifications",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		hlog.Fatal("err: channel.ExchangeDeclare : " + err.Error())
	}

	err = channel.Qos(
		1, 0,
		false,
	)
	if err != nil {
		hlog.Fatal("err: channel.Qos" + err.Error())
	}

	zap.L().Info(fmt.Sprintf("successfully connected to rabbitmq %s as consumer ", cfg.RabbitMQ.RMQAddress))

	return &RabbitMQ{
		Connection: conn,
		Channel:    channel,
	}

}

func (r *RabbitMQ) Close(ctx context.Context) {
	zap.L().Info("closing rabbitmq gracefully")
	r.Connection.Close()
}
