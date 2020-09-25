package service

import (
	"context"
	"github.com/koind/action-log/api/internal/domain/repository"
	"github.com/pkg/errors"
)

// Сервис историй действий
type HistoryService struct {
	HistoryRepository repository.HistoryRepositoryInterface
}

// Создает новый сервис
func NewHistoryService(hr repository.HistoryRepositoryInterface) *HistoryService {
	return &HistoryService{
		HistoryRepository: hr,
	}
}

// Добавляет новую историю действий
func (s *HistoryService) Add(ctx context.Context, history repository.History) (*repository.History, error) {
	history.SetDatetimeOfCreate()

	newHistory, err := s.HistoryRepository.Add(ctx, history)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка при добавлении истории действий")
	}

	return newHistory, nil
}

// Возвращает все истории действий
func (s *HistoryService) GetAll(ctx context.Context) ([]*repository.History, error) {
	list, err := s.HistoryRepository.GetAll(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка при поиске всех действий")
	}

	return list, nil
}
