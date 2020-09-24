package main

import (
	"context"
	"github.com/koind/action-log/internal/config"
	"github.com/koind/action-log/internal/domain/service"
	"github.com/koind/action-log/internal/handler"
	"github.com/koind/action-log/internal/storage/postgres"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var configPath string

const DefaultConfigPath = "config/production/config.toml"

func init() {
	pflag.StringVarP(&configPath, "config", "c", DefaultConfigPath, "Путь до конфигурационного файла")
}

func main() {
	pflag.Parse()

	cfg, err := config.Init(configPath)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(cfg.Postgres.PingTimeout)*time.Millisecond,
	)
	defer cancel()

	pg, err := cfg.IntPostgres(ctx)
	if err != nil {
		log.Fatal(err)
	}

	historyRepository := postgres.NewHistoryRepository(pg)
	historyService := service.NewHistoryService(historyRepository)
	srv := handler.NewHTTPServer(historyService, cfg.HTTPServer.GetDomain())

	// handle shutdown gracefully
	quit := make(chan os.Signal, 1)
	done := make(chan error, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		err := srv.Shutdown(ctx)

		done <- err
	}()

	logrus.Infof("Запуск сервера, %s", cfg.HTTPServer.GetDomain())
	logrus.Infof("Результат запуска сервера, %v", srv.Start())

	err = <-done
	logrus.Infof("Остановка сервера, %v", err)
}
