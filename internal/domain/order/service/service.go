package service

import (
	"context"
	"fmt"
	"time"

	"order-queue-bot/internal/domain/order"
)

type Service struct {
	repo order.Repository
}

func NewService(repo order.Repository) *Service {
	return &Service{repo: repo}
}

// CreateOrder Создание заказа
func (s *Service) CreateOrder(ctx context.Context, userID int64, description string) (*order.Order, error) {
	if description == "" {
		return nil, fmt.Errorf("description cant be empty")
	}

	newOrder := &order.Order{
		UserID:      userID,
		Description: description,
		Status:      order.StatusNew,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.Create(ctx, newOrder); err != nil {
		return nil, fmt.Errorf("failed to create order in repository: %w", err)
	}

	return newOrder, nil
}

// GetOrder Возвщарает заказ по ID

func (s *Service) GetOrder(ctx context.Context, id int64) (*order.Order, error) {
	return s.repo.GetByID(ctx, id)
}

// ChangeStatus перевод заказа из одного статуса в другой
func (s *Service) ChangeStatus(ctx context.Context, id int64, newStatus order.Status) (*order.Order, error) {
	ord, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order for status update: %w", err)
	}

	if !ord.CanTransitionTo(newStatus) {
		return nil, fmt.Errorf("%w: transition from %s to %s is forbidden",
			order.ErrInvalidStatus, ord.Status, newStatus)
	}

	ord.Status = newStatus
	ord.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, ord); err != nil {
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	return ord, nil
}
