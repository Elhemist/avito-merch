package repository

import (
	"fmt"
	"merch-test/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type InventoryPostgres struct {
	db *sqlx.DB
}

func NewInventoryPostgres(db *sqlx.DB) *InventoryPostgres {
	return &InventoryPostgres{db: db}
}

func (r *InventoryPostgres) GetItemById(itemID int) (models.MerchItem, error) {
	var merch models.MerchItem

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1;", merchTable)
	err := r.db.Get(&merch, query, itemID)
	if err != nil {
		return models.MerchItem{}, fmt.Errorf("no item with id %d found: %w", itemID, err)
	}

	return merch, err
}

func (r *InventoryPostgres) GetItemByName(itemName string) (models.MerchItem, error) {
	var merch models.MerchItem

	query := fmt.Sprintf("SELECT * FROM %s WHERE name = $1;", merchTable)
	err := r.db.Get(&merch, query, itemName)
	if err != nil {
		return models.MerchItem{}, fmt.Errorf("no item with name %s found: %w", itemName, err)
	}

	return merch, err
}

func (r *InventoryPostgres) GetUserInventory(userID uuid.UUID) ([]models.UserInventoryItem, error) {
	var Inventory []models.UserInventoryItem

	query := fmt.Sprintf(`SELECT * FROM %s WHERE user_id = $1;`, inventoryTable)
	err := r.db.Select(&Inventory, query, userID)
	if err != nil {
		return []models.UserInventoryItem{}, fmt.Errorf("inventory get error: %w", err)
	}

	return Inventory, err
}

func (r *InventoryPostgres) BuyItem(userID, walletID uuid.UUID, itemId int) error {

	tx := r.db.MustBegin()

	merch, err := r.GetItemById(itemId)
	if err != nil {
		tx.Rollback()
		return err
	}

	query := fmt.Sprintf(`UPDATE %s SET balance = balance - $1 WHERE id = $2 AND balance >= $1`, walletTable)
	result, err := tx.Exec(query, merch.Price, walletID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to spend coins: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to spend coins: %w", err)
	}
	if rowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("not enough coins")
	}

	query = fmt.Sprintf(`UPDATE %s SET quantity = quantity + 1 WHERE user_id = $1 AND merch_item_id = $2`, inventoryTable)
	result, err = tx.Exec(query, userID, itemId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update inventory: %w", err)
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		query = fmt.Sprintf(`INSERT INTO %s (user_id, merch_item_id, quantity) VALUES ($1, $2, $3)`, inventoryTable)
		_, err = tx.Exec(query, userID, itemId, 1)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to add new item to inventory: %w", err)
		}
	}

	tx.Commit()

	return nil
}
