package repository_test

import (
	"errors"
	"merch-test/internal/models"
	"merch-test/internal/repository"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestGetUserTransfersReceived(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repository.NewTransfersPostgres(sqlxDB)

	userID := uuid.New()
	transferTime := time.Now()

	expectedTransfers := []models.CoinTransfers{
		{
			ID:         1,
			SenderID:   uuid.New(),
			ReceiverID: userID,
			Amount:     100,
			CreatedAt:  transferTime,
		},
		{
			ID:         2,
			SenderID:   uuid.New(),
			ReceiverID: userID,
			Amount:     200,
			CreatedAt:  transferTime,
		},
	}

	rows := sqlmock.NewRows([]string{"id", "sender_id", "receiver_id", "amount", "created_at"})
	for _, transfer := range expectedTransfers {
		rows.AddRow(transfer.ID, transfer.SenderID, transfer.ReceiverID, transfer.Amount, transfer.CreatedAt)
	}

	{
		mock.ExpectQuery("^SELECT \\* FROM transactions WHERE receiver_id = \\$1;$").
			WithArgs(userID).
			WillReturnRows(rows)

		transfers, err := repo.GetUserTransfersReceived(userID)
		assert.NoError(t, err)
		assert.Equal(t, expectedTransfers, transfers)
	}

	{
		mock.ExpectQuery("^SELECT \\* FROM transactions WHERE receiver_id = \\$1;$").
			WithArgs(userID).
			WillReturnError(errors.New("query error"))

		_, err = repo.GetUserTransfersReceived(userID)
		assert.Error(t, err)
	}
}

func TestGetUserTransfersSent(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repository.NewTransfersPostgres(sqlxDB)

	userID := uuid.New()
	transferTime := time.Now()

	expectedTransfers := []models.CoinTransfers{
		{
			ID:         1,
			SenderID:   userID,
			ReceiverID: uuid.New(),
			Amount:     150,
			CreatedAt:  transferTime,
		},
		{
			ID:         2,
			SenderID:   userID,
			ReceiverID: uuid.New(),
			Amount:     250,
			CreatedAt:  transferTime,
		},
	}

	rows := sqlmock.NewRows([]string{"id", "sender_id", "receiver_id", "amount", "created_at"})
	for _, transfer := range expectedTransfers {
		rows.AddRow(transfer.ID, transfer.SenderID, transfer.ReceiverID, transfer.Amount, transfer.CreatedAt)
	}

	{
		mock.ExpectQuery("^SELECT \\* FROM transactions WHERE sender_id = \\$1;$").
			WithArgs(userID).
			WillReturnRows(rows)

		transfers, err := repo.GetUserTransfersSent(userID)
		assert.NoError(t, err)
		assert.Equal(t, expectedTransfers, transfers)
	}

	{
		mock.ExpectQuery("^SELECT \\* FROM transactions WHERE sender_id = \\$1;$").
			WithArgs(userID).
			WillReturnError(errors.New("query error"))

		_, err = repo.GetUserTransfersSent(userID)
		assert.Error(t, err)
	}
}
