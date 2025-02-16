package repositories

import (
	"go/avito-test/internal/models"
	"go/avito-test/pkg/db"
)

type OrderRepository struct {
	db *db.Db
}

func NewOrderRepository(db *db.Db) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}
