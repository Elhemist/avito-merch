package models

import (
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
	WalletID uuid.UUID `db:"id"`
	UserID   uuid.UUID `db:"user_id"`
	Balance  int       `db:"balance"`
}

type CoinTransfers struct {
	ID         int       `db:"id"`
	SenderID   uuid.UUID `db:"sender_id"`
	ReceiverID uuid.UUID `db:"receiver_id"`
	Amount     int       `db:"amount"`
	CreatedAt  time.Time `db:"created_at"`
}

type CoinHistory struct {
	Received []CoinTransaction `json:"received"`
	Sent     []CoinTransaction `json:"sent"`
}

type CoinTransaction struct {
	FromUser string `json:"fromUser,omitempty"`
	ToUser   string `json:"toUser,omitempty"`
	Amount   int    `json:"amount"`
}
