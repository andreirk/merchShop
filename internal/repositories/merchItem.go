package repositories

import (
	"go/avito-test/internal/models"
	"go/avito-test/pkg/db"
)

type ItemRepository struct {
	db *db.Db
}

func NewItemRepository(db *db.Db) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) FindByName(name string) (*models.MerchItem, error) {
	var item models.MerchItem
	if err := r.db.Where("name = ?", name).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
