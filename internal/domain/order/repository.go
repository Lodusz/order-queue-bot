package order

import "context"

type Repository interface {
	Create(ctx context.Context, ord *Order) error
	GetByID(ctx context.Context, id int64) (*Order, error)
	GetByUserID(ctx context.Context, userID int64) ([]Order, error)
	GetQueue(ctx context.Context) ([]Order, error)
	Update(ctx context.Context, ord *Order) error
}
