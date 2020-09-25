package postgres

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/koind/action-log/api/internal/domain/repository"
	"github.com/pkg/errors"
)

const (
	queryHealthCheck   = `SELECT 1;`
	queryInsertHistory = `INSERT INTO histories(user_id, project, description, created_at)
		VALUES ($1, $2, $3, $4) RETURNING id`
	queryFindAll = `SELECT * FROM histories`
)

// Postgres репозиторий историй действий
type HistoryRepository struct {
	DB *sqlx.DB
}

// Возвращает Postgres репозиторий историй действий
func NewHistoryRepository(db *sqlx.DB) *HistoryRepository {
	return &HistoryRepository{
		DB: db,
	}
}

// Проверяет состояние базы
func (r *HistoryRepository) HealthCheck(ctx context.Context) error {
	err := r.DB.PingContext(ctx)
	if err != nil {
		return err
	}

	var one int
	err = r.DB.QueryRowContext(ctx, queryHealthCheck).Scan(&one)

	return err
}

// Добавляет новую историю действий
func (r *HistoryRepository) Add(ctx context.Context, history repository.History) (*repository.History, error) {
	err := r.DB.QueryRowContext(
		ctx,
		queryInsertHistory,
		history.UserID,
		history.Project,
		history.Description,
		history.CreatedAt,
	).Scan(&history.ID)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка при добавлении истории действий")
	}

	return &history, nil
}

// Возвращает все истории действий
func (r *HistoryRepository) GetAll(ctx context.Context) ([]*repository.History, error) {
	rows, err := r.DB.QueryxContext(ctx, queryFindAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "ошибка при поиске историй действий")
	}

	histories := make([]*repository.History, 0)

	for rows.Next() {
		var history repository.History

		err := rows.StructScan(&history)
		if err != nil {
			return nil, errors.Wrap(err, "ошибка при сканировании результатов в структуру")
		}

		histories = append(histories, &history)
	}

	return histories, nil
}
