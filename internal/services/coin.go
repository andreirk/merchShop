package services

import (
	"database/sql"
	"errors"
	"fmt"
	"go/avito-test/internal/models"
	"go/avito-test/internal/repositories"
	logLib "log"
)

type ICoinService interface {
	SendCoins(senderName string, receiverUsername string, amount int) error
}

type CoinService struct {
	userRepo      *repositories.UserRepository
	coinTransRepo *repositories.CoinTransactionRepository
}

func NewCoinService(userRepo *repositories.UserRepository, coinTransRepo *repositories.CoinTransactionRepository) *CoinService {
	return &CoinService{userRepo: userRepo, coinTransRepo: coinTransRepo}
}

func (s *CoinService) SendCoins(senderName, receiverUsername string, amount int) error {
	const op = "internal.services.CoinService.SendCoins"
	log := logLib.New(logLib.Writer(), op, logLib.LstdFlags)
	// Start a database transaction
	db := s.userRepo.DB()
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Find sender
	sender, err := s.userRepo.FindUserByName(senderName)
	if err != nil {
		return errors.New("sender not found")
	}

	// Validate sender balance
	if sender.CoinBalance < amount {
		return errors.New("insufficient coins")
	}

	// Find receiver
	receiver, err := s.userRepo.FindUserByName(receiverUsername)
	if err != nil {
		return errors.New("receiver not found")
	}

	/// Perform transaction
	// Take from sender
	rowsAffected := tx.Exec(`
					UPDATE users
						SET coin_balance = coin_balance - @amount
					WHERE id = @sender_id and coin_balance - @amount >= 0;
	`,
		sql.Named("amount", amount), sql.Named("sender_id", sender.ID),
	).RowsAffected
	if rowsAffected == 0 {
		tx.Rollback()
		log.Printf("transaction failed in %s with updating coin_balance rows_affected %s", op, rowsAffected)
		return fmt.Errorf("transaction failed in %s with %s", op, err)
	}
	// Add to receiver
	rowsAffected = tx.Exec(`
					UPDATE users
						SET coin_balance = coin_balance + @amount
					WHERE id = @receiver_id;
`,
		sql.Named("amount", amount), sql.Named("receiver_id", receiver.ID),
	).RowsAffected

	if rowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("transaction failed in %s with %s", op, err)
	}

	// Record transaction
	transaction := &models.CoinTransaction{
		SenderID:   sender.ID,
		ReceiverID: receiver.ID,
		Amount:     amount,
	}
	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to record transaction in %s with %s", op, err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction in %s with %s", op, err)
	}
	return nil
}
