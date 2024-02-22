package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/Hymiside/lamoda-api/pkg/handler"
	"github.com/Hymiside/lamoda-api/pkg/models"
	"github.com/Hymiside/lamoda-api/pkg/repository"
	"github.com/Hymiside/lamoda-api/pkg/server"
	"github.com/Hymiside/lamoda-api/pkg/service"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := godotenv.Load(); err != nil {
		log.Panicf("error to load .env file: %v", err)
	}

	db, err := repository.NewPostgresDB(ctx, models.ConfigPostgres{
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Name:     os.Getenv("POSTGRES_DATABASE"),
	})
	if err != nil {
		log.Panicf("error to connect postgres: %v", err)
	}
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
		select {
		case <-quit:
			cancel()
		case <-ctx.Done():
			return
		}
	}()

	srv := server.Server{}
	if err = srv.RunServer(ctx, handlers.NewRoutes(), models.ConfigServer{
		Host: os.Getenv("SERVER_HOST"),
		Port: os.Getenv("SERVER_PORT")}); err != nil {
		log.Panicf("failed to run server: %v", err)
	}
}
