package postgres

import (
	"github.com/FudSy/Diploma/internal/dto"
	"github.com/FudSy/Diploma/internal/pkg/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthPostgres struct {
	db *gorm.DB
}

func NewAuthPostgres(db *gorm.DB) *AuthPostgres {
	return &AuthPostgres{db}
}

func (r *AuthPostgres) CreateUser(user dto.User) (uuid.UUID, error) {
	modelUser := models.User{
		ID:           user.ID,
		Login:        user.Login,
		Email:        user.Email,
		Name:         user.Name,
		Surname:      user.Surname,
		PasswordHash: user.PasswordHash,
		FullName:     user.FullName,
		Role:         user.Role,
	}

	if err := r.db.Create(&modelUser).Error; err != nil {
		return uuid.Nil, err
	}
	return modelUser.ID, nil
}

func (r *AuthPostgres) GetUserByLogin(login string) (dto.User, error) {
	var modelUser models.User
	err := r.db.Where("login = ?", login).First(&modelUser).Error
	if err != nil {
		return dto.User{}, err
	}

	return dto.User{
		ID:           modelUser.ID,
		Login:        modelUser.Login,
		Email:        modelUser.Email,
		Name:         modelUser.Name,
		Surname:      modelUser.Surname,
		PasswordHash: modelUser.PasswordHash,
		FullName:     modelUser.FullName,
		Role:         modelUser.Role,
	}, nil
}

func (r *AuthPostgres) GetUserById(id uuid.UUID) (dto.User, error) {
	var modelUser models.User
	err := r.db.Where("id = ?", id).First(&modelUser).Error
	if err != nil {
		return dto.User{}, err
	}

	return dto.User{
		ID:           modelUser.ID,
		Login:        modelUser.Login,
		Email:        modelUser.Email,
		Name:         modelUser.Name,
		Surname:      modelUser.Surname,
		PasswordHash: modelUser.PasswordHash,
		FullName:     modelUser.FullName,
		Role:         modelUser.Role,
	}, nil
}
