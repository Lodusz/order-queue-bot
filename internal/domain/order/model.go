package order

import (
	"errors"
	"time"
)

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
	StatusCancelled  Status = "cancelled"
)

// errOrderNotfound - oshibka kogda zakaz ne naiden

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrInvalidStatus = errors.New("invalid order status")
)

type Order struct {
	ID          int64
	UserID      int64 // ID usera v telegrame
	Description string
	Status      Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (o *Order) CanTransitionTo(newStatus Status) bool {
	switch o.Status {
	case StatusNew:
		return newStatus == StatusInProgress || newStatus == StatusCancelled
	case StatusInProgress:
		return newStatus == StatusDone || newStatus == StatusCancelled
	case StatusDone, StatusCancelled:
		return false
	default:
		return false
	}
}
