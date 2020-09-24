package service

import (
	"context"
	"github.com/koind/action-log/internal/domain/repository"
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
func (s *HistoryService) Add(ctx context.Context, rotation repository.History) (*repository.History, error) {
	newRotation, err := s.HistoryRepository.Add(ctx, rotation)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка при добавлении истории действий")
	}

	return newRotation, nil
}

// Ищет истории действий по фильтрам
func (s *HistoryService) FindAllByFilter(ctx context.Context, filter repository.SearchFilter) ([]*repository.History, error) {
	list, err := s.HistoryRepository.FindAllByFilter(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка при поиске историй действий")
	}

	return list, nil
}
