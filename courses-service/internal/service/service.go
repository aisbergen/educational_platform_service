package service

import (
	"context"
	"courses/internal/domain"
	"courses/internal/repository"
	"courses/pkg/cache"
	"courses/pkg/kafka"
	"time"

	_ "github.com/golang/mock/mockgen/model"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go
type Courses interface {
	Create(ctx context.Context, course domain.Courses) error
	GetByID(ctx context.Context, id int) (domain.Courses, error)
	Update(ctx context.Context, course domain.Courses) error
	Delete(ctx context.Context, id int) error
	GetCoursesByIdStudent(ctx context.Context, studentId string) ([]domain.Courses, error)
}

type Kafka interface {
	SendMessages(topic string, message string) error
	ConsumeMessages(topic string, handler func(message string)) error
	// Read(ctx context.Context)
	ConsumeResponseMessages()
	Close()
}

type Service struct {
	Courses Courses
	Kafka   *KafkaService
}

func NewService(repo *repository.Repository, cache cache.Cache, ttl time.Duration, producer *kafka.Producer, concumer *kafka.Consumer) *Service {
	return &Service{
		Courses: NewCoursesService(repo.Courses, cache, ttl),
		Kafka:   NewKafkaSerivce(producer, concumer, repo.Courses),
	}
}
