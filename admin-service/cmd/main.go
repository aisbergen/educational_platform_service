package main

import (
	"admin/internal/config"
	delivery "admin/internal/delivery/http"
	"admin/internal/repository"
	"admin/internal/server"
	"admin/internal/service"
	"admin/pkg/auth"
	"admin/pkg/database"
	"admin/pkg/hash"
	"admin/pkg/kafka"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	path = "./.env"
)

// @title Admin Service API
// @version 1.0
// @description API Server for Admin Application

// @host localhost:8080
// @BasePath /api/v1/

// @securityDefinitions.apikey AdminAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.NewConfig(path)
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.OpenDB(cfg.Database.Driver, cfg.Database.DSN)

	if err != nil {
		log.Fatalf("error creating database object: %v", err)
	}

	hasher := hash.NewHash(cfg.Hash.Cost)

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

	service := service.NewService(repos, hasher, tokenManager, cfg, producer, consumer)

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
