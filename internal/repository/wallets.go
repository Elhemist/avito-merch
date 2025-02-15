package repository

import (
	"fmt"
	"merch-test/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type WalletPostgres struct {
	db *sqlx.DB
}

func NewWalletPostgres(db *sqlx.DB) *WalletPostgres {
	return &WalletPostgres{db: db}
}

func (r *WalletPostgres) GetUserWallet(userID uuid.UUID) (models.Wallet, error) {
	var wallet models.Wallet

	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1;", walletTable)
	err := r.db.Get(&wallet, query, userID)
	if err != nil {
		return models.Wallet{}, fmt.Errorf("failed to get wallet for user %s: %w", userID, err)
	}
	return wallet, err
}

func (r *WalletPostgres) CreateTransaction(senderID, receiverID uuid.UUID, amount int) error {
	senderWallet, err := r.GetUserWallet(senderID)
	if err != nil {
		return err
	}

	receiverWallet, err := r.GetUserWallet(receiverID)
	if err != nil {
		return err
	}

	if senderWallet.Balance < amount {
		return fmt.Errorf("insufficient balance")
	}

	tx := r.db.MustBegin()
	query := fmt.Sprintf(`UPDATE %s SET balance = balance - $1 WHERE id = $2 AND balance >= $1`, walletTable)
	res, err := tx.Exec(query, amount, senderWallet.WalletID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to send coins: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return fmt.Errorf("insufficient balance")
	}

	query = fmt.Sprintf(`UPDATE %s SET balance = balance + $1 WHERE id = $2`, walletTable)
	_, err = tx.Exec(query, amount, receiverWallet.WalletID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to add coins: %w", err)
	}

	query = fmt.Sprintf(`INSERT INTO %s (sender_id, receiver_id, amount) VALUES ($1, $2, $3)`, transferTable)
	_, err = tx.Exec(query, senderID, receiverID, amount)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to record transaction: %w", err)
	}

	tx.Commit()

	return nil
}
