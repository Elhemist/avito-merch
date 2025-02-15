package service

import (
	"merch-test/internal/models"
	"merch-test/internal/repository"

	"github.com/google/uuid"
)

type Authorization interface {
	CreateUser(user models.AuthRequest) (uuid.UUID, error)
	GenerateToken(user models.AuthRequest) (string, error)
	ParseToken(token string) (uuid.UUID, error)
}
type Inventory interface {
	BuyItem(userID uuid.UUID, purchase string) error
}
type Wallet interface {
	SendCoin(senderID uuid.UUID, receiverName string, amount int) error
}

type User interface {
	GetInfo(userID uuid.UUID) (models.InfoResponse, error)
}
type Service struct {
	Authorization
	Inventory
	Wallet
	User
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.UserRepository),
		Inventory:     NewInventoryService(repos.InventoryRepository, repos.WalletRepository),
		Wallet:        NewWalletService(repos.WalletRepository, repos.UserRepository),
		User:          NewUserService(repos.WalletRepository, repos.InventoryRepository, repos.TransfersRepository, repos.UserRepository),
	}
}
