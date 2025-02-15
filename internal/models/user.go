package models

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID `db:"id"`
	Username     string    `db:"username"`
	PasswordHash string    `db:"password_hash"`
}

type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type InfoResponse struct {
	Coins       int            `json:"coins"`
	Inventory   []UserPurchase `json:"inventory"`
	CoinHistory CoinHistory    `json:"coinHistory"`
}
