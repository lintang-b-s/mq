package dal

import (
	"lintang/go_consumer/biz/dal/rabbitmq"
	"lintang/go_consumer/config"
)

func InitRmq(cfg *config.Config) *rabbitmq.RabbitMQ {
	rmq := rabbitmq.NewRabbitMQ(cfg)

	return rmq
}
