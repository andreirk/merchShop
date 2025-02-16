package repositories

import (
	"go/avito-test/internal/models"
	"go/avito-test/pkg/db"
)

type UserRepository struct {
	db *db.Db
}

func NewUserRepository(db *db.Db) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindUserByName(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Preload("Orders.Item").
		Preload("SentCoinTransactions.Receiver").
		Preload("ReceivedCoinTransactions.Sender").
		Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(user *models.User) (*models.User, error) {
	result := r.db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}
func (r *UserRepository) FindOrCreate(user *models.User) (*models.User, error) {
	result := r.db.FirstOrCreate(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Preload(relation string) *UserRepository {
	r.db.Preload(relation)

	return r
}

func (r *UserRepository) First(user *models.User, id uint) error {
	return r.db.First(user, id).Error
}

func (r *UserRepository) DB() *db.Db {
	return r.db
}
