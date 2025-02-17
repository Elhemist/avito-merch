package repository

import (
	"fmt"
	"merch-test/internal/models"

	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) GetUserById(userID uuid.UUID) (models.User, error) {
	var user models.User

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1;", userTable)
	err := r.db.Get(&user, query, userID)
	if err != nil {
		return models.User{}, fmt.Errorf("no user with id: %d found: %w", userID, err)
	}
	return user, err
}

func (r *UserPostgres) GetUserByName(username string) (models.User, error) {
	var user models.User

	query := fmt.Sprintf("SELECT * FROM %s WHERE username = $1;", userTable)
	err := r.db.Get(&user, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, nil
		}
		return models.User{}, fmt.Errorf("no user with name: %s found: %w", username, err)
	}
	return user, err
}

func (r *UserPostgres) CreateUser(user models.AuthRequest, balance int) (uuid.UUID, error) {
	var userID uuid.UUID
	tx := r.db.MustBegin()
	query := fmt.Sprintf(`INSERT INTO %s (username, password_hash) VALUES ($1, $2) RETURNING id;`, userTable)
	err := tx.QueryRow(query, user.Username, user.Password).Scan(&userID)
	if err != nil {
		tx.Rollback()
		return uuid.Nil, fmt.Errorf("user create error: %w", err)
	}
	query = fmt.Sprintf(`INSERT INTO %s (user_id, balance) VALUES ($1, $2)`, walletTable)
	_, err = tx.Exec(query, userID, balance)
	if err != nil {
		tx.Rollback()
		return uuid.Nil, fmt.Errorf("wallet create error: %w", err)
	}
	tx.Commit()
	return userID, err
}
