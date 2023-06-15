package service

import (
	"admin/internal/config"
	"admin/internal/domain"
	"admin/internal/repository"
	"admin/pkg/auth"
	"admin/pkg/hash"
	"admin/pkg/kafka"
	"context"
)

type Admin interface {
	SignUp(ctx context.Context, admin domain.Admin) error
	SignIn(ctx context.Context, email string, password string) (domain.Token, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (domain.Token, error)
}

type Kafka interface {
	SendMessages(topic string, message string) error
	ConsumeMessages(topic string, handler func(message string)) error
}

type Service struct {
	Admin *AdminService
	Kafka Kafka
}

func NewService(repo *repository.Repository, hash hash.PasswordHasher, tokenManager auth.TokenManager, cfg *config.Config, producer *kafka.Producer, consumer *kafka.Consumer) *Service {
	return &Service{
		Admin: NewAdminService(repo.Admin, hash, tokenManager, cfg.JWT.AccessTokenTTL, cfg.JWT.RefreshTokenTTL),
		Kafka: NewKafkaSerivce(producer, consumer),
	}
}
