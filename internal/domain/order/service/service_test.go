package service

import (
	"context"
	"testing"

	"order-queue-bot/internal/domain/order"
)

// zaglushka
type mockRepo struct {
	orders map[int64]*order.Order
	idSeq  int64
}

func newMockRepo() *mockRepo {
	return &mockRepo{orders: make(map[int64]*order.Order)}
}

func (m *mockRepo) Create(ctx context.Context, ord *order.Order) error {
	m.idSeq++
	ord.ID = m.idSeq
	m.orders[ord.ID] = ord
	return nil
}

func (m *mockRepo) GetByID(ctx context.Context, id int64) (*order.Order, error) {
	ord, ok := m.orders[id]
	if !ok {
		return nil, order.ErrOrderNotFound
	}
	copyOrd := *ord
	return &copyOrd, nil
}

func (m *mockRepo) GetByUserID(ctx context.Context, userID int64) ([]order.Order, error) {
	return nil, nil
}
func (m *mockRepo) GetQueue(ctx context.Context) ([]order.Order, error) { return nil, nil }

func (m *mockRepo) Update(ctx context.Context, ord *order.Order) error {
	if _, ok := m.orders[ord.ID]; !ok {
		return order.ErrOrderNotFound
	}
	m.orders[ord.ID] = ord
	return nil
}

func TestChangeStatus(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)
	ctx := context.Background()

	_ = repo.Create(ctx, &order.Order{
		UserID: 1, Description: "Test", Status: order.StatusNew,
	})

	tests := []struct {
		name      string
		orderID   int64
		newStatus order.Status
		wantErr   bool
	}{
		{"valid: new -> in_progress", 1, order.StatusInProgress, false},
		{"invalid: in_progress -> new", 1, order.StatusNew, true},
		{"not found", 999, order.StatusDone, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.ChangeStatus(ctx, tt.orderID, tt.newStatus)
			if (err != nil) != tt.wantErr {
				t.Errorf("wantErr %v, got error %v", tt.wantErr, err)
			}
		})
	}
}
