package dal

import (
	"lintang/go_producer/biz/dal/rabbitmq"
	"lintang/go_producer/config"
)

func InitRmq(cfg *config.Config) *rabbitmq.RabbitMQ {
	rmq := rabbitmq.NewRabbitMQ(cfg)

	return rmq
}
