package repository

import (
	"merch-test/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	CreateUser(user models.AuthRequest, balance int) (uuid.UUID, error)
	GetUserByName(username string) (models.User, error)
	GetUserById(username uuid.UUID) (models.User, error)
}
type TransfersRepository interface {
	GetUserTransfersSent(userID uuid.UUID) ([]models.CoinTransfers, error)
	GetUserTransfersReceived(userID uuid.UUID) ([]models.CoinTransfers, error)
}

type WalletRepository interface {
	GetUserWallet(userId uuid.UUID) (models.Wallet, error)
	CreateTransaction(senderWallet, receiverWallet uuid.UUID, amount int) error
}

type InventoryRepository interface {
	GetUserInventory(userID uuid.UUID) ([]models.UserInventoryItem, error)
	GetItemById(itemID int) (models.MerchItem, error)
	GetItemByName(itemName string) (models.MerchItem, error)
	BuyItem(userID, walletID uuid.UUID, itemId int) error
}

type Repository struct {
	InventoryRepository
	WalletRepository
	UserRepository
	TransfersRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UserRepository:      NewUserPostgres(db),
		WalletRepository:    NewWalletPostgres(db),
		InventoryRepository: NewInventoryPostgres(db),
		TransfersRepository: NewTransfersPostgres(db),
	}
}
