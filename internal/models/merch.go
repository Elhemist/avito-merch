package models

import "github.com/google/uuid"

type MerchItem struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Price int    `db:"price"`
}

type UserInventoryItem struct {
	ID          int       `db:"id"`
	UserID      uuid.UUID `db:"user_id"`
	MerchItemID int       `db:"merch_item_id"`
	Quantity    int       `db:"quantity"`
}

type UserPurchase struct {
	Name     string `db:"merch_name"`
	Quantity int    `db:"quantity"`
}
