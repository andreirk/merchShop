package services

import (
	"database/sql"
	"errors"
	"fmt"
	"go/avito-test/internal/models"
	"go/avito-test/internal/repositories"
	logLib "log"
)

type IOrderService interface {
	PurchaseItem(username, itemName string) error
}

type OrderService struct {
	userRepo  *repositories.UserRepository
	itemRepo  *repositories.ItemRepository
	orderRepo *repositories.OrderRepository
}

func NewOrderService(userRepo *repositories.UserRepository, itemRepo *repositories.ItemRepository, orderRepo *repositories.OrderRepository) *OrderService {
	return &OrderService{userRepo: userRepo, itemRepo: itemRepo, orderRepo: orderRepo}
}

func (s *OrderService) PurchaseItem(username, itemName string) error {
	const op = "internal.services.OrderService.PurchaseItem"
	log := logLib.New(logLib.Writer(), op, logLib.LstdFlags)
	db := s.userRepo.DB()
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// Find user
	user, err := s.userRepo.FindUserByName(username)
	if err != nil {
		return errors.New("user not found")
	}

	// Find item
	item, err := s.itemRepo.FindByName(itemName)
	if err != nil {
		return errors.New("item not found")
	}

	// Validate user balance
	if user.CoinBalance < item.Price {
		return errors.New("insufficient coins")
	}

	/// Perform transaction
	// Deduct coins
	rowsAffected := tx.Exec(`
			UPDATE users
				SET coin_balance = coin_balance - @item_price
			WHERE id = @user_id and coin_balance - @item_price >= 0;
	`,
		sql.Named("item_price", item.Price), sql.Named("user_id", user.ID),
	).RowsAffected
	if rowsAffected == 0 {
		tx.Rollback()
		log.Printf("transaction failed in %s with updating coin_balance rows_affected %s", op, rowsAffected)
		return errors.New("transaction failed")
	}

	// Record order
	order := &models.Order{
		UserID:   user.ID,
		ItemID:   item.ID,
		Quantity: 1,
	}
	if err := tx.Create(order).Error; err != nil {
		errMsg := fmt.Sprintf("failed to create order in %s: with %s", op, err)
		log.Println(errMsg)
		return fmt.Errorf(errMsg)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		errMsg := fmt.Sprintf("failed to commit transaction in %s: with %s", op, err)
		log.Println(errMsg)
		return fmt.Errorf(errMsg)
	}

	return nil
}
