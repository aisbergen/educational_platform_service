package service

import (
	"context"
	"student/internal/config"
	"student/internal/domain"
	"student/internal/repository"
	"student/pkg/auth"
	"student/pkg/cache"
	"student/pkg/hash"
	"student/pkg/kafka"
)

//go:generate mockgen -source=services.go -destination=mocks/mock.go
type Students interface {
	Create(ctx context.Context, student domain.Student) error
	GetStudentByID(ctx context.Context, id int) (domain.Student, error)
	Update(ctx context.Context, student domain.Student) error
	Delete(ctx context.Context, id int) error
	GetStudentsByCoursesID(ctx context.Context, id string) ([]domain.Student, error)
	GetByEmail(ctx context.Context, email string, password string) (domain.Token, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (domain.Token, error)
}

type Kafka interface {
	Read(ctx context.Context)
	SendMessages(topic string, message string) error
	ConsumeMessages(topic string, handler func(message string)) error
	Close()
}

type Service struct {
	Students Students
	Kafka    Kafka
}

func NewService(repo *repository.Repository, hash hash.PasswordHasher, tokenManager auth.TokenManager, cache cache.Cache, cfg *config.Config, producer *kafka.Producer, consumer *kafka.Consumer) *Service {
	return &Service{
		Students: NewStudentService(repo.Students, hash, tokenManager, cache, cfg.JWT.AccessTokenTTL, cfg.JWT.RefreshTokenTTL),
		Kafka:    NewKafkaSerivce(repo.Students, producer, consumer),
	}
}
