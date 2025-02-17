package repository_test

import (
	"errors"
	"merch-test/internal/models"
	"merch-test/internal/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestGetUserWallet(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repository.NewWalletPostgres(sqlxDB)

	userID := uuid.New()
	expectedWallet := models.Wallet{
		WalletID: uuid.New(),
		UserID:   userID,
		Balance:  500,
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "balance"}).
		AddRow(expectedWallet.WalletID, expectedWallet.UserID, expectedWallet.Balance)

	{
		mock.ExpectQuery("^SELECT \\* FROM wallets WHERE user_id = \\$1;$").
			WithArgs(userID).
			WillReturnRows(rows)

		wallet, err := repo.GetUserWallet(userID)
		assert.NoError(t, err)
		assert.Equal(t, expectedWallet, wallet)
	}

	{
		mock.ExpectQuery("^SELECT \\* FROM wallets WHERE user_id = \\$1;$").
			WithArgs(uuid.New()).
			WillReturnError(errors.New("failed to get wallet for user"))

		_, err = repo.GetUserWallet(uuid.New())
		assert.Error(t, err)
	}
}

func TestCreateTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repository.NewWalletPostgres(sqlxDB)

	senderID := uuid.New()
	receiverID := uuid.New()
	amount := 200

	senderWallet := models.Wallet{
		WalletID: uuid.New(),
		UserID:   senderID,
		Balance:  500,
	}

	receiverWallet := models.Wallet{
		WalletID: uuid.New(),
		UserID:   receiverID,
		Balance:  300,
	}

	mock.ExpectQuery("^SELECT \\* FROM wallets WHERE user_id = \\$1;$").
		WithArgs(senderID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance"}).AddRow(senderWallet.WalletID, senderWallet.UserID, senderWallet.Balance))

	mock.ExpectQuery("^SELECT \\* FROM wallets WHERE user_id = \\$1;$").
		WithArgs(receiverID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance"}).AddRow(receiverWallet.WalletID, receiverWallet.UserID, receiverWallet.Balance))

	{
		mock.ExpectBegin()
		mock.ExpectExec("^UPDATE wallets SET balance = balance - \\$1 WHERE id = \\$2 AND balance >= \\$1$").
			WithArgs(amount, senderWallet.WalletID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec("^UPDATE wallets SET balance = balance \\+ \\$1 WHERE id = \\$2$").
			WithArgs(amount, receiverWallet.WalletID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec("^INSERT INTO transactions \\(sender_id, receiver_id, amount\\) VALUES \\(\\$1, \\$2, \\$3\\)$").
			WithArgs(senderID, receiverID, amount).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err = repo.CreateTransaction(senderID, receiverID, amount)
		assert.NoError(t, err)
	}
	{
		senderWallet.Balance = 100
		mock.ExpectQuery("^SELECT \\* FROM wallets WHERE user_id = \\$1;$").
			WithArgs(senderID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance"}).AddRow(senderWallet.WalletID, senderWallet.UserID, senderWallet.Balance))

		mock.ExpectQuery("^SELECT \\* FROM wallets WHERE user_id = \\$1;$").
			WithArgs(receiverID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance"}).AddRow(receiverWallet.WalletID, receiverWallet.UserID, receiverWallet.Balance))

		err = repo.CreateTransaction(senderID, receiverID, amount)
		assert.Error(t, err)
		assert.Equal(t, "insufficient balance", err.Error())
	}
	{
		mock.ExpectBegin()
		mock.ExpectExec("^UPDATE wallets SET balance = balance - \\$1 WHERE id = \\$2 AND balance >= \\$1$").
			WithArgs(amount, senderWallet.WalletID).
			WillReturnError(errors.New("failed to send coins"))
		mock.ExpectRollback()

		err = repo.CreateTransaction(senderID, receiverID, amount)
		assert.Error(t, err)
	}
	{
		mock.ExpectBegin()
		mock.ExpectExec("^UPDATE wallets SET balance = balance - \\$1 WHERE id = \\$2 AND balance >= \\$1$").
			WithArgs(amount, senderWallet.WalletID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("^UPDATE wallets SET balance = balance \\+ \\$1 WHERE id = \\$2$").
			WithArgs(amount, receiverWallet.WalletID).
			WillReturnError(errors.New("failed to add coins"))
		mock.ExpectRollback()

		err = repo.CreateTransaction(senderID, receiverID, amount)
		assert.Error(t, err)

	}
	{
		mock.ExpectBegin()
		mock.ExpectExec("^UPDATE wallets SET balance = balance - \\$1 WHERE id = \\$2 AND balance >= \\$1$").
			WithArgs(amount, senderWallet.WalletID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("^UPDATE wallets SET balance = balance \\+ \\$1 WHERE id = \\$2$").
			WithArgs(amount, receiverWallet.WalletID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("^INSERT INTO transactions \\(sender_id, receiver_id, amount\\) VALUES \\(\\$1, \\$2, \\$3\\)$").
			WithArgs(senderID, receiverID, amount).
			WillReturnError(errors.New("failed to record transaction"))
		mock.ExpectRollback()

		err = repo.CreateTransaction(senderID, receiverID, amount)
		assert.Error(t, err)
	}
}
