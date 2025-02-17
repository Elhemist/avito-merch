package service_test

import (
	"merch-test/internal/models"
	"merch-test/internal/service"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserRepo struct {
	mock.Mock
}

const (
	salt       = "someSalt"
	signingKey = "podpis"
	tokenTTL   = time.Hour / 2
)

func (m *mockUserRepo) CreateUser(user models.AuthRequest, balance int) (uuid.UUID, error) {
	args := m.Called(user, balance)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *mockUserRepo) GetUserByName(username string) (models.User, error) {
	args := m.Called(username)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *mockUserRepo) GetUserById(userID uuid.UUID) (models.User, error) {
	args := m.Called(userID)
	return args.Get(0).(models.User), args.Error(1)
}

func TestCreateUser(t *testing.T) {
	mockRepo := new(mockUserRepo)
	authService := service.NewAuthService(mockRepo)

	user := models.AuthRequest{
		Username: "testUser",
		Password: "password123",
	}
	userID := uuid.New()

	mockRepo.On("CreateUser", mock.Anything, service.NEW_USER_BALANCE).Return(userID, nil)

	id, err := authService.CreateUser(user)
	assert.NoError(t, err)
	assert.Equal(t, userID, id)

	mockRepo.AssertExpectations(t)
}

func TestGenerateToken_ExistingUser(t *testing.T) {
	mockRepo := new(mockUserRepo)
	authService := service.NewAuthService(mockRepo)

	userReq := models.AuthRequest{
		Username: "existingUser",
		Password: "password123",
	}
	userID := uuid.New()
	hashedPassword := service.GeneratePasswordHash(userReq.Password)

	mockRepo.On("GetUserByName", userReq.Username).Return(models.User{ID: userID, PasswordHash: hashedPassword}, nil)

	token, err := authService.GenerateToken(userReq)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	mockRepo.AssertExpectations(t)
}

func TestGenerateToken_NewUser(t *testing.T) {
	mockRepo := new(mockUserRepo)
	authService := service.NewAuthService(mockRepo)

	userReq := models.AuthRequest{
		Username: "newUser",
		Password: "password123",
	}
	userID := uuid.New()

	mockRepo.On("GetUserByName", userReq.Username).Return(models.User{}, nil)
	mockRepo.On("CreateUser", mock.Anything, service.NEW_USER_BALANCE).Return(userID, nil)

	token, err := authService.GenerateToken(userReq)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	mockRepo.AssertExpectations(t)
}

func TestGenerateToken_WrongPassword(t *testing.T) {
	mockRepo := new(mockUserRepo)
	authService := service.NewAuthService(mockRepo)

	userReq := models.AuthRequest{
		Username: "existingUser",
		Password: "wrongPassword",
	}
	userID := uuid.New()
	correctPasswordHash := service.GeneratePasswordHash("correctPassword")

	mockRepo.On("GetUserByName", userReq.Username).Return(models.User{ID: userID, PasswordHash: correctPasswordHash}, nil)

	token, err := authService.GenerateToken(userReq)
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "Unauthorized", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestParseToken_Success(t *testing.T) {
	authService := service.NewAuthService(nil)

	userID := uuid.New()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &service.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId: userID,
	})
	tokenString, _ := token.SignedString([]byte(signingKey))

	parsedUserID, err := authService.ParseToken(tokenString)
	assert.NoError(t, err)
	assert.Equal(t, userID, parsedUserID)
}

func TestParseToken_InvalidSignature(t *testing.T) {
	authService := service.NewAuthService(nil)

	userID := uuid.New()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &service.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId: userID,
	})
	tokenString, _ := token.SignedString([]byte("wrongSigningKey"))

	parsedUserID, err := authService.ParseToken(tokenString)
	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, parsedUserID)
}

func TestParseToken_InvalidToken(t *testing.T) {
	authService := service.NewAuthService(nil)

	parsedUserID, err := authService.ParseToken("invalid.token.string")
	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, parsedUserID)
}
