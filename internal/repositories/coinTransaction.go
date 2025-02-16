package repositories

import (
	"go/avito-test/internal/models"
	"go/avito-test/pkg/db"
)

type CoinTransactionRepository struct {
	db *db.Db
}

func NewCoinTransactionRepository(db *db.Db) *CoinTransactionRepository {
	return &CoinTransactionRepository{db: db}
}

func (r *CoinTransactionRepository) Create(transaction *models.CoinTransaction) error {
	return r.db.Create(transaction).Error
}
