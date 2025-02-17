package service_test

import (
	"errors"
	"merch-test/internal/models"
	"merch-test/internal/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockInventoryRepo struct {
	mock.Mock
}

func (m *mockInventoryRepo) GetUserInventory(userID uuid.UUID) ([]models.UserInventoryItem, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.UserInventoryItem), args.Error(1)
}

func (m *mockInventoryRepo) GetItemById(itemID int) (models.MerchItem, error) {
	args := m.Called(itemID)
	return args.Get(0).(models.MerchItem), args.Error(1)
}

func (m *mockInventoryRepo) GetItemByName(itemName string) (models.MerchItem, error) {
	args := m.Called(itemName)
	return args.Get(0).(models.MerchItem), args.Error(1)
}

func (m *mockInventoryRepo) BuyItem(userID, walletID uuid.UUID, itemId int) error {
	args := m.Called(userID, walletID, itemId)
	return args.Error(0)
}

func TestBuyItem_Success(t *testing.T) {
	mockInventory := new(mockInventoryRepo)
	mockWallet := new(mockWalletRepo)
	service := service.NewInventoryService(mockInventory, mockWallet)

	userID := uuid.New()
	walletID := uuid.New()
	itemID := 1
	itemName := "Test Item"
	item := models.MerchItem{ID: itemID, Name: itemName}
	wallet := models.Wallet{WalletID: walletID}

	mockInventory.On("GetItemByName", itemName).Return(item, nil)
	mockWallet.On("GetUserWallet", userID).Return(wallet, nil)
	mockInventory.On("BuyItem", userID, walletID, itemID).Return(nil)

	err := service.BuyItem(userID, itemName)

	assert.NoError(t, err)
	mockInventory.AssertExpectations(t)
	mockWallet.AssertExpectations(t)
}

func TestBuyItem_ItemNotFound(t *testing.T) {
	mockInventory := new(mockInventoryRepo)
	mockWallet := new(mockWalletRepo)
	service := service.NewInventoryService(mockInventory, mockWallet)

	userID := uuid.New()
	itemName := "Nonexistent Item"

	mockInventory.On("GetItemByName", itemName).Return(models.MerchItem{}, errors.New("item not found"))

	err := service.BuyItem(userID, itemName)

	assert.Error(t, err)
	assert.Equal(t, "item not found", err.Error())
	mockInventory.AssertExpectations(t)
}

func TestBuyItem_WalletError(t *testing.T) {
	mockInventory := new(mockInventoryRepo)
	mockWallet := new(mockWalletRepo)
	service := service.NewInventoryService(mockInventory, mockWallet)

	userID := uuid.New()
	itemID := 1
	itemName := "Test Item"
	item := models.MerchItem{ID: itemID, Name: itemName}

	mockInventory.On("GetItemByName", itemName).Return(item, nil)
	mockWallet.On("GetUserWallet", userID).Return(models.Wallet{}, errors.New("failed to get wallet"))

	err := service.BuyItem(userID, itemName)

	assert.Error(t, err)
	assert.Equal(t, "failed to get wallet", err.Error())
	mockInventory.AssertExpectations(t)
	mockWallet.AssertExpectations(t)
}

func TestBuyItem_BuyError(t *testing.T) {
	mockInventory := new(mockInventoryRepo)
	mockWallet := new(mockWalletRepo)
	service := service.NewInventoryService(mockInventory, mockWallet)

	userID := uuid.New()
	walletID := uuid.New()
	itemID := 1
	itemName := "Test Item"
	item := models.MerchItem{ID: itemID, Name: itemName}
	wallet := models.Wallet{WalletID: walletID}

	mockInventory.On("GetItemByName", itemName).Return(item, nil)
	mockWallet.On("GetUserWallet", userID).Return(wallet, nil)
	mockInventory.On("BuyItem", userID, walletID, itemID).Return(errors.New("purchase error"))

	err := service.BuyItem(userID, itemName)

	assert.Error(t, err)
	assert.Equal(t, "purchase error", err.Error())
	mockInventory.AssertExpectations(t)
	mockWallet.AssertExpectations(t)
}
