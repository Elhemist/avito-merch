package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"merch-test/internal/models"
	"merch-test/internal/repository"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const NEW_USER_BALANCE = 1000
const (
	salt       = "someSalt"
	signingKey = "podpis"
	tokenTTL   = time.Hour / 2
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId uuid.UUID `json:"user_id"`
}

type AuthorizationService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthorizationService {
	return &AuthorizationService{userRepo: userRepo}
}

func (s *AuthorizationService) CreateUser(user models.AuthRequest) (uuid.UUID, error) {
	user.Password = generatePasswordHash(user.Password)

	return s.userRepo.CreateUser(user, NEW_USER_BALANCE)
}

func (s *AuthorizationService) GenerateToken(userReq models.AuthRequest) (string, error) {
	user, err := s.userRepo.GetUserByName(userReq.Username)
	if err != nil {
		logrus.Info(err)
		return "", err
	}
	if user == (models.User{}) {
		logrus.Info("creating new user:", userReq.Username)
		id, err := s.CreateUser(userReq)
		if err != nil {
			logrus.Info(err)
			return "", err
		}
		user.ID = id
	} else if user.PasswordHash != generatePasswordHash(userReq.Password) {
		return "", fmt.Errorf("Unauthorized")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.ID,
	})

	return token.SignedString([]byte(signingKey))
}

func (s *AuthorizationService) ParseToken(accessToken string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return uuid.Nil, errors.New("invalid token struct")
	}

	return claims.UserId, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
