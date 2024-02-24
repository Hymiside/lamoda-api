package main

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"runtime"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/Hymiside/lamoda-api/pkg/handler"
	"github.com/Hymiside/lamoda-api/pkg/models"
	"github.com/Hymiside/lamoda-api/pkg/repository"
	"github.com/Hymiside/lamoda-api/pkg/service"
	log "github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
            _, filename := path.Split(f.File)
            filename = fmt.Sprintf("%s:%d", filename, f.Line)
            return "", filename
        },
	})

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
		log.Fatalf("error to connect postgres: %v", err)
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

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT")),
		Handler: handlers.NewRoutes(),
	}

	go func() {
		<-ctx.Done()
		if err := httpServer.Shutdown(context.Background()); err != nil {
			log.Fatalf("failed to shutdown server: %v", err)
		}
	}()

	log.Infof("server started on http://%s:%s/", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))
	if err = httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to run server: %v", err)
	}
}
