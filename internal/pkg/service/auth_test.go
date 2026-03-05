package service

import (
	"errors"
	"os"
	"testing"

	"github.com/FudSy/Diploma/internal/dto"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authRepoMock struct {
	createUserFn     func(user dto.User) (uuid.UUID, error)
	getUserByLoginFn func(login string) (dto.User, error)
	getUserByIDFn    func(id uuid.UUID) (dto.User, error)
}

func (m *authRepoMock) CreateUser(user dto.User) (uuid.UUID, error) {
	return m.createUserFn(user)
}

func (m *authRepoMock) GetUserByLogin(login string) (dto.User, error) {
	return m.getUserByLoginFn(login)
}

func (m *authRepoMock) GetUserById(id uuid.UUID) (dto.User, error) {
	return m.getUserByIDFn(id)
}

func TestAuthService_CreateUser_Duplicate(t *testing.T) {
	repo := &authRepoMock{
		createUserFn: func(user dto.User) (uuid.UUID, error) { return uuid.Nil, nil },
		getUserByLoginFn: func(login string) (dto.User, error) {
			return dto.User{ID: uuid.New(), Login: login}, nil
		},
		getUserByIDFn: func(id uuid.UUID) (dto.User, error) { return dto.User{}, nil },
	}
	svc := NewAuthService(repo)

	_, err := svc.CreateUser(dto.RegisterRequest{
		Login: "user1", Email: "u1@example.com", Name: "Ivan", Surname: "Petrov", Password: "secret12",
	})
	if err == nil {
		t.Fatal("expected duplicate error, got nil")
	}
}

func TestAuthService_CreateAdmin_SetsRole(t *testing.T) {
	var captured dto.User
	repo := &authRepoMock{
		createUserFn: func(user dto.User) (uuid.UUID, error) {
			captured = user
			return uuid.New(), nil
		},
		getUserByLoginFn: func(login string) (dto.User, error) { return dto.User{}, errors.New("not found") },
		getUserByIDFn:    func(id uuid.UUID) (dto.User, error) { return dto.User{}, nil },
	}
	svc := NewAuthService(repo)

	_, err := svc.CreateAdmin(dto.RegisterRequest{
		Login: "admin1", Email: "a1@example.com", Name: "Anna", Surname: "Smirnova", Password: "secret12",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if captured.Role != "ADMIN" {
		t.Fatalf("expected ADMIN role, got %q", captured.Role)
	}
	if captured.PasswordHash == "secret12" {
		t.Fatal("password was not hashed")
	}
}

func TestAuthService_Login_SuccessAndParseToken(t *testing.T) {
	_ = os.Setenv("AUTH_SIGNING_KEY", "test-signing-key")
	defer os.Unsetenv("AUTH_SIGNING_KEY")

	userID := uuid.New()
	hash, err := bcrypt.GenerateFromPassword([]byte("secret12"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash error: %v", err)
	}
	repo := &authRepoMock{
		createUserFn: func(user dto.User) (uuid.UUID, error) { return uuid.New(), nil },
		getUserByLoginFn: func(login string) (dto.User, error) {
			return dto.User{ID: userID, Login: login, PasswordHash: string(hash)}, nil
		},
		getUserByIDFn: func(id uuid.UUID) (dto.User, error) { return dto.User{ID: id}, nil },
	}
	svc := NewAuthService(repo)

	token, err := svc.Login("user1", "secret12")
	if err != nil {
		t.Fatalf("unexpected login error: %v", err)
	}

	parsed, err := svc.ParseToken(token)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if parsed != userID {
		t.Fatalf("expected user id %s, got %s", userID, parsed)
	}
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("secret12"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash error: %v", err)
	}
	repo := &authRepoMock{
		createUserFn: func(user dto.User) (uuid.UUID, error) { return uuid.New(), nil },
		getUserByLoginFn: func(login string) (dto.User, error) {
			return dto.User{ID: uuid.New(), Login: login, PasswordHash: string(hash)}, nil
		},
		getUserByIDFn: func(id uuid.UUID) (dto.User, error) { return dto.User{ID: id}, nil },
	}
	svc := NewAuthService(repo)

	_, err = svc.Login("user1", "wrong-pass")
	if err == nil {
		t.Fatal("expected invalid credentials error, got nil")
	}
}

func TestAuthService_ParseToken_InvalidFormat(t *testing.T) {
	svc := NewAuthService(&authRepoMock{
		createUserFn:     func(user dto.User) (uuid.UUID, error) { return uuid.New(), nil },
		getUserByLoginFn: func(login string) (dto.User, error) { return dto.User{}, errors.New("not found") },
		getUserByIDFn:    func(id uuid.UUID) (dto.User, error) { return dto.User{}, nil },
	})

	_, err := svc.ParseToken("not-base64")
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
}
