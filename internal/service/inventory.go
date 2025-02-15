package service

import (
	"merch-test/internal/repository"

	"github.com/google/uuid"
)

type InventoryService struct {
	inventoryRepo repository.InventoryRepository
	walletRepo    repository.WalletRepository
}

func NewInventoryService(inventoryRepo repository.InventoryRepository, walletRepo repository.WalletRepository) *InventoryService {
	return &InventoryService{inventoryRepo: inventoryRepo, walletRepo: walletRepo}
}

func (s *InventoryService) BuyItem(userID uuid.UUID, purchase string) error {
	item, err := s.inventoryRepo.GetItemByName(purchase)
	if err != nil {
		return err
	}

	wallet, err := s.walletRepo.GetUserWallet(userID)
	if err != nil {
		return err
	}

	err = s.inventoryRepo.BuyItem(userID, wallet.WalletID, item.ID)
	if err != nil {
		return err
	}

	return nil
}
