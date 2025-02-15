package service

import (
	"merch-test/internal/repository"

	"github.com/google/uuid"
)

type WalletService struct {
	walletRepo repository.WalletRepository
	userRepo   repository.UserRepository
}

func NewWalletService(WalletRepo repository.WalletRepository, UserRepo repository.UserRepository) *WalletService {
	return &WalletService{walletRepo: WalletRepo, userRepo: UserRepo}
}

func (s *WalletService) SendCoin(senderID uuid.UUID, receiverName string, amount int) error {

	receiver, err := s.userRepo.GetUserByName(receiverName)
	if err != nil {
		return err
	}

	err = s.walletRepo.CreateTransaction(senderID, receiver.ID, amount)

	return err
}
