package repository_test

import (
	"database/sql"
	"errors"
	"merch-test/internal/models"
	"merch-test/internal/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestGetUserById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repository.NewUserPostgres(sqlxDB)

	userID := uuid.New()
	expectedUser := models.User{
		ID:           userID,
		Username:     "testuser",
		PasswordHash: "hashedpassword",
	}

	rows := sqlmock.NewRows([]string{"id", "username", "password_hash"}).
		AddRow(expectedUser.ID, expectedUser.Username, expectedUser.PasswordHash)

	{
		mock.ExpectQuery("^SELECT \\* FROM users WHERE id = \\$1;$").
			WithArgs(userID).
			WillReturnRows(rows)

		user, err := repo.GetUserById(userID)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	}

	{
		mock.ExpectQuery("^SELECT \\* FROM users WHERE id = \\$1;$").
			WithArgs(uuid.New()).
			WillReturnError(sql.ErrNoRows)

		_, err = repo.GetUserById(uuid.New())
		assert.Error(t, err)
	}
}

func TestGetUserByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repository.NewUserPostgres(sqlxDB)

	username := "testuser"
	expectedUser := models.User{
		ID:           uuid.New(),
		Username:     username,
		PasswordHash: "hashedpassword",
	}

	rows := sqlmock.NewRows([]string{"id", "username", "password_hash"}).
		AddRow(expectedUser.ID, expectedUser.Username, expectedUser.PasswordHash)

	{
		mock.ExpectQuery("^SELECT \\* FROM users WHERE username = \\$1;$").
			WithArgs(username).
			WillReturnRows(rows)

		user, err := repo.GetUserByName(username)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	}
	{
		mock.ExpectQuery("^SELECT \\* FROM users WHERE username = \\$1;$").
			WithArgs("ramdom").
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUserByName("ramdom")
		assert.NoError(t, err)
		assert.Equal(t, models.User{}, user)
	}
}

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repository.NewUserPostgres(sqlxDB)

	userRequest := models.AuthRequest{
		Username: "testuser",
		Password: "testpassword",
	}
	userID := uuid.New()
	balance := 100

	{
		mock.ExpectBegin()
		mock.ExpectQuery("^INSERT INTO users \\(username, password_hash\\) VALUES \\(\\$1, \\$2\\) RETURNING id;$").
			WithArgs(userRequest.Username, userRequest.Password).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))

		mock.ExpectExec("^INSERT INTO wallets \\(user_id, balance\\) VALUES \\(\\$1, \\$2\\)").
			WithArgs(userID, balance).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		id, err := repo.CreateUser(userRequest, balance)
		assert.NoError(t, err)
		assert.Equal(t, userID, id)
	}
	{
		mock.ExpectBegin()
		mock.ExpectQuery("^INSERT INTO users \\(username, password_hash\\) VALUES \\(\\$1, \\$2\\) RETURNING id;$").
			WithArgs(userRequest.Username, userRequest.Password).
			WillReturnError(errors.New("insert user error"))
		mock.ExpectRollback()

		_, err = repo.CreateUser(userRequest, balance)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user create error")
	}
}
