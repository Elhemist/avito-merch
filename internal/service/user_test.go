package service_test

import (
	"errors"
	"merch-test/internal/models"
	"merch-test/internal/service"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockTransfersRepo struct {
	mock.Mock
}

func (m *mockTransfersRepo) GetUserTransfersReceived(userID uuid.UUID) ([]models.CoinTransfers, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.CoinTransfers), args.Error(1)
}

func (m *mockTransfersRepo) GetUserTransfersSent(userID uuid.UUID) ([]models.CoinTransfers, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.CoinTransfers), args.Error(1)
}

func TestGetInfo_Success(t *testing.T) {
	mockWallet := new(mockWalletRepo)
	mockInventory := new(mockInventoryRepo)
	mockTransfers := new(mockTransfersRepo)
	mockUser := new(mockUserRepo)
	userService := service.NewUserService(mockWallet, mockInventory, mockTransfers, mockUser)

	userID := uuid.New()
	wallet := models.Wallet{Balance: 100}
	inventory := []models.UserInventoryItem{
		{ID: 1, UserID: userID, MerchItemID: 1, Quantity: 2},
	}
	merchItem := models.MerchItem{Name: "pink-hoody"}
	receivedTransfers := []models.CoinTransfers{
		{ID: 1, SenderID: uuid.New(), ReceiverID: userID, Amount: 50, CreatedAt: time.Now()},
	}
	sentTransfers := []models.CoinTransfers{
		{ID: 2, SenderID: userID, ReceiverID: uuid.New(), Amount: 20, CreatedAt: time.Now()},
	}
	userSender := models.User{ID: receivedTransfers[0].SenderID, Username: "Alice"}
	userReceiver := models.User{ID: sentTransfers[0].ReceiverID, Username: "Bob"}

	mockWallet.On("GetUserWallet", userID).Return(wallet, nil)
	mockInventory.On("GetUserInventory", userID).Return(inventory, nil)
	mockInventory.On("GetItemById", inventory[0].MerchItemID).Return(merchItem, nil)
	mockTransfers.On("GetUserTransfersReceived", userID).Return(receivedTransfers, nil)
	mockTransfers.On("GetUserTransfersSent", userID).Return(sentTransfers, nil)
	mockUser.On("GetUserById", receivedTransfers[0].SenderID).Return(userSender, nil)
	mockUser.On("GetUserById", sentTransfers[0].ReceiverID).Return(userReceiver, nil)

	info, err := userService.GetInfo(userID)
	assert.NoError(t, err)
	assert.Equal(t, 100, info.Coins)
	assert.Len(t, info.Inventory, 1)
	assert.Equal(t, "pink-hoody", info.Inventory[0].Name)
	assert.Len(t, info.CoinHistory.Received, 1)
	assert.Equal(t, "Alice", info.CoinHistory.Received[0].FromUser)
	assert.Len(t, info.CoinHistory.Sent, 1)
	assert.Equal(t, "Bob", info.CoinHistory.Sent[0].ToUser)

	mockWallet.AssertExpectations(t)
	mockInventory.AssertExpectations(t)
	mockTransfers.AssertExpectations(t)
	mockUser.AssertExpectations(t)
}
func TestGetInfo_WalletError(t *testing.T) {
	mockWallet := new(mockWalletRepo)
	mockInventory := new(mockInventoryRepo)
	mockTransfers := new(mockTransfersRepo)
	mockUser := new(mockUserRepo)
	userService := service.NewUserService(mockWallet, mockInventory, mockTransfers, mockUser)

	userID := uuid.New()
	mockWallet.On("GetUserWallet", userID).Return(models.Wallet{}, errors.New("wallet error"))

	info, err := userService.GetInfo(userID)
	assert.Error(t, err)
	assert.Equal(t, "wallet error", err.Error())
	assert.Equal(t, service.NewEmptyInfoResponse(), info)

	mockWallet.AssertExpectations(t)
}

func TestGetInfo_InventoryError(t *testing.T) {
	mockWallet := new(mockWalletRepo)
	mockInventory := new(mockInventoryRepo)
	mockTransfers := new(mockTransfersRepo)
	mockUser := new(mockUserRepo)
	userService := service.NewUserService(mockWallet, mockInventory, mockTransfers, mockUser)

	userID := uuid.New()
	wallet := models.Wallet{Balance: 100}

	mockWallet.On("GetUserWallet", userID).Return(wallet, nil)
	mockInventory.On("GetUserInventory", userID).Return([]models.UserInventoryItem{}, errors.New("inventory error"))

	info, err := userService.GetInfo(userID)

	assert.Error(t, err)
	assert.Equal(t, service.NewEmptyInfoResponse(), info)

	mockWallet.AssertExpectations(t)
	mockInventory.AssertExpectations(t)
}

func TestGetInfo_TransferError(t *testing.T) {
	mockWallet := new(mockWalletRepo)
	mockInventory := new(mockInventoryRepo)
	mockTransfers := new(mockTransfersRepo)
	mockUser := new(mockUserRepo)
	userService := service.NewUserService(mockWallet, mockInventory, mockTransfers, mockUser)

	userID := uuid.New()
	wallet := models.Wallet{Balance: 100}

	mockWallet.On("GetUserWallet", userID).Return(wallet, nil)
	mockInventory.On("GetUserInventory", userID).Return([]models.UserInventoryItem{}, nil)
	mockTransfers.On("GetUserTransfersReceived", userID).Return([]models.CoinTransfers{}, errors.New("transfer error"))

	info, err := userService.GetInfo(userID)

	assert.Error(t, err)
	assert.Equal(t, service.NewEmptyInfoResponse(), info)

	mockWallet.AssertExpectations(t)
	mockInventory.AssertExpectations(t)
	mockTransfers.AssertExpectations(t)
}
