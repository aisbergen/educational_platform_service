package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/subosito/gotenv"
	"golang.org/x/crypto/bcrypt"
)

const (
	defaultServerPort               = "8080"
	defaultServerRWTimeout          = 10 * time.Second
	defaultServerMaxHeaderMegabytes = 1
	defaultAccessTokenTTL           = 15 * time.Minute
	defaultRefreshTokenTTL          = 24 * time.Hour * 30
	defailtCost                     = bcrypt.DefaultCost
)

type JWTConfig struct {
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	SigningKey      string
}

type serverConfig struct {
	Port               string
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	MaxHeaderMegabytes int
}

type databaseConfig struct {
	Driver string
	DSN    string
}

type hashConfig struct {
	Cost int
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
}

type Config struct {
	JWT      JWTConfig
	Server   serverConfig
	Database databaseConfig
	Hash     hashConfig
	Kafka    KafkaConfig
}

func NewConfig(path string) (*Config, error) {
	err := gotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load environment variables from file: %v", err)
	}

	driver := os.Getenv("DRIVER")
	dsn := os.Getenv("DSN_STUDENTS")
	jwtKey := os.Getenv("SIGNIN_KEY")

	brokerStr := os.Getenv("KAFKA_BROKERS")
	brokers := strings.Split(brokerStr, ",")
	topic := os.Getenv("KAFKA_TOPIC")

	return &Config{
		JWT: JWTConfig{
			AccessTokenTTL:  defaultAccessTokenTTL,
			RefreshTokenTTL: defaultRefreshTokenTTL,
			SigningKey:      jwtKey,
		},
		Server: serverConfig{
			Port:               defaultServerPort,
			ReadTimeout:        defaultServerRWTimeout,
			WriteTimeout:       defaultServerRWTimeout,
			MaxHeaderMegabytes: defaultServerMaxHeaderMegabytes,
		},
		Database: databaseConfig{
			Driver: driver,
			DSN:    dsn,
		},
		Hash: hashConfig{
			Cost: defailtCost,
		},
		Kafka: KafkaConfig{
			Brokers: brokers,
			Topic:   topic,
		},
	}, nil
}
