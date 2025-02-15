package service

import (
	"fmt"
	"merch-test/internal/models"
	"merch-test/internal/repository"

	"github.com/google/uuid"
)

type UserService struct {
	walletRepo    repository.WalletRepository
	inventoryRepo repository.InventoryRepository
	transferRepo  repository.TransfersRepository
	userRepo      repository.UserRepository
}

func NewUserService(walletRepo repository.WalletRepository, inventoryRepo repository.InventoryRepository, transferRepo repository.TransfersRepository, userRepo repository.UserRepository) *UserService {
	return &UserService{walletRepo: walletRepo, inventoryRepo: inventoryRepo, transferRepo: transferRepo, userRepo: userRepo}
}
func NewEmptyInfoResponse() models.InfoResponse {
	return models.InfoResponse{
		Coins:     0,
		Inventory: []models.UserPurchase{},
		CoinHistory: models.CoinHistory{
			Received: []models.CoinTransaction{},
			Sent:     []models.CoinTransaction{},
		},
	}
}

func (s *UserService) convertInventory(inventory []models.UserInventoryItem) ([]models.UserPurchase, error) {
	purchases := make([]models.UserPurchase, 0, len(inventory))

	for _, item := range inventory {
		merchItem, err := s.inventoryRepo.GetItemById(item.MerchItemID)
		if err != nil {
			return nil, fmt.Errorf("failed to get merch item for ID %d: %w", item.MerchItemID, err)
		}

		purchase := models.UserPurchase{
			Name:     merchItem.Name,
			Quantity: item.Quantity,
		}

		purchases = append(purchases, purchase)
	}

	return purchases, nil
}

func (s *UserService) convertReceivedTransfers(transfers []models.CoinTransfers) ([]models.CoinTransaction, error) {
	userIDs := make(map[uuid.UUID]struct{})
	for _, transfer := range transfers {
		userIDs[transfer.SenderID] = struct{}{}
	}
	userList := make(map[uuid.UUID]string)
	for id := range userIDs {
		user, err := s.userRepo.GetUserById(id)
		if err != nil {
			return nil, err
		}
		userList[id] = user.Username
	}
	transactions := make([]models.CoinTransaction, len(transfers))
	for i, transfer := range transfers {
		transactions[i] = models.CoinTransaction{
			FromUser: userList[transfer.SenderID],
			Amount:   transfer.Amount,
		}
	}
	return transactions, nil
}

func (s *UserService) convertSentTransfers(transfers []models.CoinTransfers) ([]models.CoinTransaction, error) {
	userIDs := make(map[uuid.UUID]struct{})
	for _, transfer := range transfers {
		userIDs[transfer.ReceiverID] = struct{}{}
	}
	userList := make(map[uuid.UUID]string)
	for id := range userIDs {
		user, err := s.userRepo.GetUserById(id)
		if err != nil {
			return nil, err
		}
		userList[id] = user.Username
	}
	transactions := make([]models.CoinTransaction, len(transfers))
	for i, transfer := range transfers {
		transactions[i] = models.CoinTransaction{
			ToUser: userList[transfer.ReceiverID],
			Amount: transfer.Amount,
		}
	}
	return transactions, nil
}

func (s *UserService) GetInfo(userID uuid.UUID) (models.InfoResponse, error) {

	wallet, err := s.walletRepo.GetUserWallet(userID)
	if err != nil {
		return NewEmptyInfoResponse(), err
	}

	inventory, err := s.inventoryRepo.GetUserInventory(userID)
	if err != nil {
		return NewEmptyInfoResponse(), err
	}

	purchase, err := s.convertInventory(inventory)
	if err != nil {
		return NewEmptyInfoResponse(), err
	}

	receivedTransfers, err := s.transferRepo.GetUserTransfersReceived(userID)
	if err != nil {
		return NewEmptyInfoResponse(), err
	}
	receivedConv, err := s.convertReceivedTransfers(receivedTransfers)
	if err != nil {
		return NewEmptyInfoResponse(), err
	}

	sentTransfers, err := s.transferRepo.GetUserTransfersSent(userID)
	if err != nil {
		return NewEmptyInfoResponse(), err
	}
	sentConv, err := s.convertSentTransfers(sentTransfers)
	if err != nil {
		return NewEmptyInfoResponse(), err
	}

	infoResponse := models.InfoResponse{
		Coins:     wallet.Balance,
		Inventory: purchase,
		CoinHistory: models.CoinHistory{
			Received: receivedConv,
			Sent:     sentConv,
		},
	}
	return infoResponse, err
}
