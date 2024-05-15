package config

import (
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App
		HTTP
		// Redis
		LogConfig

		RabbitMQ

		Kafka
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	LogConfig struct {
		Level string `json:"level" yaml:"level" env:"LOG_LEVEL"`
		// Filename   string `json:"filename" yaml:"filename"`
		// MaxSize    int    `json:"maxsize" yaml:"maxsize"`
		MaxAge     int `json:"max_age" yaml:"max_age" env:"LOG_MAXAGE"`
		MaxBackups int `json:"max_backups" yaml:"max_backups" env:"LOG_MAXBACKUP"`
	}

	RabbitMQ struct {
		RMQAddress string `json:"rabbitmqAddress" yaml:"rmqAddress" env:"RABBITMQ_ADDRESS"`
	}

	Kafka struct {
		KafkaAddress      string `json:"kafka_address" env:"KAFKA_ADDRESS"`
		NotificationTopic string `json:"kafka_topic" env:"NOTIFICATION_TOPIC"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	// err = cleanenv.ReadConfig(path+".env", cfg) // buat di doker , ../.env kalo debug (.env kalo docker)
	// err = cleanenv.ReadConfig(path+"/local.env", cfg) // local run
	err = cleanenv.ReadConfig(path+"/local.env", cfg)

	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
