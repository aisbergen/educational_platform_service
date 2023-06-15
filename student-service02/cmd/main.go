package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"student/internal/config"
	delivery "student/internal/delivery/http"
	"student/internal/repository"
	"student/internal/server"
	"student/internal/service"
	"student/pkg/auth"
	"student/pkg/cache"
	"student/pkg/database"
	"student/pkg/hash"
	"student/pkg/kafka"
	"time"
)

const (
	path = "./.env"
)

// @title Student Service API
// @version 1.0
// @description API Server for Student Application

// @host localhost:8001
// @BasePath /api/v1/

// @securityDefinitions.apikey StudentAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.Init(path)
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.OpenDB(cfg.Database.Driver, cfg.Database.DSN)

	if err != nil {
		log.Fatalf("error creating database object: %v", err)
	}

	hasher := hash.NewHash(cfg.Hash.Cost)

	memCache, err := cache.NewMemoryCache(cfg.Redis)

	if err != nil {
		log.Fatalf("error creating mem cache: %v", err)
	}

	tokenManager, err := auth.NewManager(cfg.JWT.SigningKey)
	if err != nil {
		log.Fatal(err)
	}

	producer, err := kafka.NewProducer(cfg.Kafka.Brokers)
	if err != nil {
		log.Fatalf("error creating Kafka producer: %v", err)
	}

	consumer, err := kafka.NewConsumer(cfg.Kafka.Brokers)
	if err != nil {
		log.Fatalf("error creating Kafka consumer: %v", err)
	}

	repos := repository.NewRepository(db)

	service := service.NewService(repos, hasher, tokenManager, memCache, cfg, producer, consumer)

	go service.Kafka.Read(context.Background())

	handler := delivery.NewHandler(service, tokenManager)

	srv := server.NewServer(cfg, handler.Init(cfg))

	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	log.Println("Server started", cfg.Server.Port)

	quit := make(chan os.Signal, 1)

	<-quit

	const timeout = 5 * time.Second
	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		log.Fatalf("failed to stop server: %v", err)
	}
}
