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

// CreateOrder создать заказ
func (s *Service) CreateOrder(ctx context.Context, userID int64, description string) (*order.Order, error) {
	if description == "" {
		return nil, fmt.Errorf("description cannot be empty")
	}

	newOrder := &order.Order{
		UserID:      userID,
		Description: description,
		Status:      order.StatusNew,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.Create(ctx, newOrder); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return newOrder, nil
}

func (s *Service) GetQueue(ctx context.Context) ([]order.Order, error) {
	return s.repo.GetQueue(ctx)
}

// GetById получаать заказ по id
func (s *Service) GetByID(ctx context.Context, id int64) (*order.Order, error) {
	return s.repo.GetByID(ctx, id)
}

// ChangeStatus поменять статус заказа ( БЕЗОПАСНО )
func (s *Service) ChangeStatus(ctx context.Context, id int64, newStatus order.Status) (*order.Order, error) {
	ord, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}

	if !ord.CanTransitionTo(newStatus) {
		return nil, fmt.Errorf("invalid status transition")
	}

	ord.Status = newStatus
	ord.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, ord); err != nil {
		return nil, fmt.Errorf("failed to update status: %w", err)
	}

	return ord, nil
}
