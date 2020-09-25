package config

import (
	"context"
	"fmt"
	"github.com/BurntSushi/toml"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
)

// Настройки микросервиса
type Options struct {
	Postgres   Postgres
	HTTPServer HTTPServer
}

// Инициализация конфигов
func Init(configPath string) (options *Options, err error) {
	if _, err = toml.DecodeFile(configPath, &options); err != nil {
		return nil, errors.Wrap(err, "не удалось загрузить конфиги микросервиса")
	}

	return
}

// Создает пул соединений и возвращает само соединение
func (o Options) IntPostgres(ctx context.Context) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", o.Postgres.DSN)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка при создании пула соединений")
	}

	db.SetMaxOpenConns(o.Postgres.MaxOpenConns)
	db.SetMaxIdleConns(o.Postgres.MaxIdleConns)
	db.SetConnMaxLifetime(o.Postgres.ConnMaxLifetime)

	if err := db.PingContext(ctx); err != nil {
		return nil, errors.Wrap(err, "ошибка при подключении к базе данных")
	}

	return db, nil
}

// Настройки базы
type Postgres struct {
	DSN             string
	PingTimeout     int
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// Настройки HTTP сервера
type HTTPServer struct {
	Host string
	Port int
}

// Возвращает домен сервера
func (s HTTPServer) GetDomain() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
