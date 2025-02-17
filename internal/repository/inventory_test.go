package repository_test

import (
	"merch-test/internal/models"
	"merch-test/internal/repository"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestGetItemById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewInventoryPostgres(sqlxDB)

	item := models.MerchItem{ID: 1, Name: "T-Shirt", Price: 100}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM merch WHERE id = $1;")).
		WithArgs(item.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price"}).
			AddRow(item.ID, item.Name, item.Price))

	result, err := repo.GetItemById(item.ID)
	assert.NoError(t, err)
	assert.Equal(t, item, result)
}

func TestGetItemByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewInventoryPostgres(sqlxDB)

	item := models.MerchItem{ID: 2, Name: "Mug", Price: 50}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM merch WHERE name = $1;")).
		WithArgs(item.Name).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price"}).
			AddRow(item.ID, item.Name, item.Price))

	result, err := repo.GetItemByName(item.Name)
	assert.NoError(t, err)
	assert.Equal(t, item, result)
}

func TestGetUserInventory(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewInventoryPostgres(sqlxDB)

	userID := uuid.New()
	inventory := []models.UserInventoryItem{
		{ID: 1, UserID: userID, MerchItemID: 1, Quantity: 2},
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM inventory WHERE user_id = $1;")).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "merch_item_id", "quantity"}).
			AddRow(inventory[0].ID, inventory[0].UserID, inventory[0].MerchItemID, inventory[0].Quantity))

	result, err := repo.GetUserInventory(userID)
	assert.NoError(t, err)
	assert.Equal(t, inventory, result)
}

func TestBuyItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewInventoryPostgres(sqlxDB)

	userID := uuid.New()
	walletID := uuid.New()
	item := models.MerchItem{ID: 1, Name: "T-Shirt", Price: 100}

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM merch WHERE id = $1;")).
		WithArgs(item.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price"}).
			AddRow(item.ID, item.Name, item.Price))

	mock.ExpectExec(regexp.QuoteMeta("UPDATE wallets SET balance = balance - $1 WHERE id = $2 AND balance >= $1")).
		WithArgs(item.Price, walletID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(regexp.QuoteMeta("UPDATE inventory SET quantity = quantity + 1 WHERE user_id = $1 AND merch_item_id = $2")).
		WithArgs(userID, item.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err = repo.BuyItem(userID, walletID, item.ID)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
