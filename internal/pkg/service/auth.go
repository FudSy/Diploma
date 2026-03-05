package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/FudSy/Diploma/internal/dto"
	"github.com/FudSy/Diploma/internal/pkg/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const tokenTTL = 12 * time.Hour

type AuthService struct {
	repo       repository.Authorization
	signingKey []byte
}

func NewAuthService(repo repository.Authorization) *AuthService {
	key := os.Getenv("AUTH_SIGNING_KEY")
	if key == "" {
		key = "local-dev-signing-key"
	}

	return &AuthService{
		repo:       repo,
		signingKey: []byte(key),
	}
}

func (s *AuthService) CreateUser(input dto.RegisterRequest) (uuid.UUID, error) {
	return s.createWithRole(input, "USER")
}

func (s *AuthService) CreateAdmin(input dto.RegisterRequest) (uuid.UUID, error) {
	return s.createWithRole(input, "ADMIN")
}

func (s *AuthService) createWithRole(input dto.RegisterRequest, role string) (uuid.UUID, error) {
	_, err := s.repo.GetUserByLogin(input.Login)
	if err == nil {
		return uuid.Nil, errors.New("user with this login already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, err
	}

	user := dto.User{
		Login:        input.Login,
		PasswordHash: string(hash),
		Email:        input.Email,
		Name:         input.Name,
		Surname:      input.Surname,
		FullName:     strings.TrimSpace(input.Name + " " + input.Surname),
		Role:         role,
	}

	return s.repo.CreateUser(user)
}

func (s *AuthService) Login(login, password string) (string, error) {
	user, err := s.repo.GetUserByLogin(login)
	if err != nil {
		return "", errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	return s.generateToken(user.ID)
}

func (s *AuthService) ParseToken(accessToken string) (uuid.UUID, error) {
	rawToken, err := base64.RawURLEncoding.DecodeString(accessToken)
	if err != nil {
		return uuid.Nil, errors.New("invalid token format")
	}

	parts := strings.Split(string(rawToken), "|")
	if len(parts) != 3 {
		return uuid.Nil, errors.New("invalid token payload")
	}

	payload := strings.Join(parts[:2], "|")
	signature := s.sign(payload)
	if !hmac.Equal([]byte(signature), []byte(parts[2])) {
		return uuid.Nil, errors.New("invalid token signature")
	}

	expiresAtUnix, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return uuid.Nil, errors.New("invalid token expiration")
	}
	if time.Now().Unix() > expiresAtUnix {
		return uuid.Nil, errors.New("token expired")
	}

	userID, err := uuid.Parse(parts[0])
	if err != nil {
		return uuid.Nil, errors.New("invalid user id in token")
	}

	return userID, nil
}

func (s *AuthService) GetUserById(id uuid.UUID) (dto.User, error) {
	return s.repo.GetUserById(id)
}

func (s *AuthService) CheckPassword(user dto.User, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
}

func (s *AuthService) generateToken(userID uuid.UUID) (string, error) {
	expiresAt := time.Now().Add(tokenTTL).Unix()
	payload := fmt.Sprintf("%s|%d", userID.String(), expiresAt)
	signature := s.sign(payload)
	token := fmt.Sprintf("%s|%s", payload, signature)

	return base64.RawURLEncoding.EncodeToString([]byte(token)), nil
}

func (s *AuthService) sign(payload string) string {
	mac := hmac.New(sha256.New, s.signingKey)
	mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
