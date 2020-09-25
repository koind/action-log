package repository

import (
	"context"
	"time"
)

// Интерфейс репозитория историй действий
type HistoryRepositoryInterface interface {
	// Проверяет состояние базы
	HealthCheck(ctx context.Context) error

	// Добавляет новую историю действий
	Add(ctx context.Context, history History) (*History, error)

	// Возвращает все истории действий
	GetAll(ctx context.Context) ([]*History, error)
}

// Сущность истории действий
type History struct {
	ID          int       `json:"id" db:"id"`                                       // Id истории
	UserID      int       `json:"userId" db:"user_id" validate:"required"`          // Id пользователя
	Project     string    `json:"project" db:"project" validate:"required"`         // Название проекта
	Description string    `json:"description" db:"description" validate:"required"` // Описание действия
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`                        // Дата создания истории
}

// Установить дату и время создания
func (r *History) SetDatetimeOfCreate() {
	r.CreatedAt = time.Now().UTC()
}
