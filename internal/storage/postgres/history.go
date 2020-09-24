package postgres

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/koind/action-log/internal/domain/repository"
	"github.com/pkg/errors"
)

const (
	queryInsertHistory = `INSERT INTO histories(user_id, project, description, created_at)
		VALUES ($1, $2, $3, $4) RETURNING id`
	queryFindAllByFilter = `SELECT * FROM histories WHERE user_id=$1 AND project=$2`
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

// Ищет истории действий по фильтрам
func (r *HistoryRepository) FindAllByFilter(ctx context.Context, filter repository.SearchFilter) ([]*repository.History, error) {
	rows, err := r.DB.QueryxContext(ctx, queryFindAllByFilter, filter.UserID, filter.Project)
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
