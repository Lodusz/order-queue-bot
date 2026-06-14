package postgres

import (
	"context"
	"errors"

	"order-queue-bot/internal/domain/order"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

type orderModel struct {
	ID          int64  `gorm:"primaryKey;autoIncrement"`
	UserID      int64  `gorm:"index"`
	Description string `gorm:"not null"`
	Status      string `gorm:"index;not null"`
	CreatedAt   int64  `gorm:"autoCreateTime"`
	UpdatedAt   int64  `gorm:"autoUpdateTime"`
}

func (orderModel) TableName() string {
	return "orders"
}

func toModel(ord *order.Order) *orderModel {
	return &orderModel{
		ID:          ord.ID,
		UserID:      ord.UserID,
		Description: ord.Description,
		Status:      string(ord.Status),
	}
}

func toDomain(m *orderModel) *order.Order {
	return &order.Order{
		ID:          m.ID,
		UserID:      m.UserID,
		Description: m.Description,
		Status:      order.Status(m.Status),
	}
}

func (r *Repository) Create(ctx context.Context, ord *order.Order) error {
	m := toModel(ord)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	ord.ID = m.ID
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*order.Order, error) {
	var m orderModel
	err := r.db.WithContext(ctx).First(&m, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, order.ErrOrderNotFound
	}
	if err != nil {
		return nil, err
	}
	return toDomain(&m), nil
}

func (r *Repository) GetQueue(ctx context.Context) ([]order.Order, error) {
	var models []orderModel
	err := r.db.WithContext(ctx).
		Where("status IN ?", []order.Status{order.StatusNew, order.StatusInProgress}).
		Order("created_at asc").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	orders := make([]order.Order, len(models))
	for i, m := range models {
		orders[i] = *toDomain(&m)
	}
	return orders, nil
}

func (r *Repository) GetByUserID(ctx context.Context, userID int64) ([]order.Order, error) {
	var models []orderModel
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	orders := make([]order.Order, len(models))
	for i, m := range models {
		orders[i] = *toDomain(&m)
	}
	return orders, nil
}

func (r *Repository) Update(ctx context.Context, ord *order.Order) error {
	m := toModel(ord)
	res := r.db.WithContext(ctx).Model(m).Select("Status", "UpdatedAt").Updates(m)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return order.ErrOrderNotFound
	}
	return nil
}
