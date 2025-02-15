package repository

import (
	"fmt"
	"merch-test/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TransfersPostgres struct {
	db *sqlx.DB
}

func NewTransfersPostgres(db *sqlx.DB) *TransfersPostgres {
	return &TransfersPostgres{db: db}
}

func (r *TransfersPostgres) GetUserTransfersReceived(userID uuid.UUID) ([]models.CoinTransfers, error) {
	var transferList []models.CoinTransfers

	query := fmt.Sprintf("SELECT * FROM %s WHERE receiver_id = $1;", transferTable)
	err := r.db.Select(&transferList, query, userID)
	return transferList, err
}

func (r *TransfersPostgres) GetUserTransfersSent(userID uuid.UUID) ([]models.CoinTransfers, error) {
	var transferList []models.CoinTransfers

	query := fmt.Sprintf("SELECT * FROM %s WHERE sender_id = $1;", transferTable)
	err := r.db.Select(&transferList, query, userID)
	return transferList, err
}
