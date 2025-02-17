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

type mockWalletRepo struct {
	mock.Mock
}

func (m *mockWalletRepo) GetUserWallet(userID uuid.UUID) (models.Wallet, error) {
	args := m.Called(userID)
	return args.Get(0).(models.Wallet), args.Error(1)
}

func (m *mockWalletRepo) CreateTransaction(senderWallet, receiverWallet uuid.UUID, amount int) error {
	args := m.Called(senderWallet, receiverWallet, amount)
	return args.Error(0)
}

func TestSendCoin_Success(t *testing.T) {
	mockWallet := new(mockWalletRepo)
	mockUser := new(mockUserRepo)
	service := service.NewWalletService(mockWallet, mockUser)

	senderID := uuid.New()
	receiverID := uuid.New()
	receiverName := "TestReceiver"
	amount := 100

	mockUser.On("GetUserByName", receiverName).Return(models.User{ID: receiverID}, nil)
	mockWallet.On("CreateTransaction", senderID, receiverID, amount).Return(nil)

	err := service.SendCoin(senderID, receiverName, amount)

	assert.NoError(t, err)
	mockUser.AssertExpectations(t)
	mockWallet.AssertExpectations(t)
}

func TestSendCoin_UserNotFound(t *testing.T) {
	mockWallet := new(mockWalletRepo)
	mockUser := new(mockUserRepo)
	service := service.NewWalletService(mockWallet, mockUser)

	senderID := uuid.New()
	receiverName := "UnknownUser"
	amount := 100

	mockUser.On("GetUserByName", receiverName).Return(models.User{}, errors.New("user not found"))

	err := service.SendCoin(senderID, receiverName, amount)

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	mockUser.AssertExpectations(t)
}

func TestSendCoin_TransactionFailed(t *testing.T) {
	mockWallet := new(mockWalletRepo)
	mockUser := new(mockUserRepo)
	service := service.NewWalletService(mockWallet, mockUser)

	senderID := uuid.New()
	receiverID := uuid.New()
	receiverName := "TestReceiver"
	amount := 100

	mockUser.On("GetUserByName", receiverName).Return(models.User{ID: receiverID}, nil)
	mockWallet.On("CreateTransaction", senderID, receiverID, amount).Return(errors.New("transaction failed"))

	err := service.SendCoin(senderID, receiverName, amount)

	assert.Error(t, err)
	assert.Equal(t, "transaction failed", err.Error())
	mockUser.AssertExpectations(t)
	mockWallet.AssertExpectations(t)
}
