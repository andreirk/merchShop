package di

import (
	"go/avito-test/internal/models"
)

type IStatRepository interface {
	AddClick(linkId uint)
}

type IUserRepository interface {
	Create(user *models.User) (*models.User, error)
	FindUserByName(username string) (*models.User, error)
}

type ICoinsTransactionRepository interface {
	Create(user *models.User) (*models.User, error)
	GetFrom(userName string) ([]*models.CoinTransaction, error)
	GetTo(userName string) ([]*models.CoinTransaction, error)
}
