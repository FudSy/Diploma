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

	return toDTOUser(modelUser), nil
}

func (r *AuthPostgres) GetUserById(id uuid.UUID) (dto.User, error) {
	var modelUser models.User
	err := r.db.Where("id = ?", id).First(&modelUser).Error
	if err != nil {
		return dto.User{}, err
	}

	return toDTOUser(modelUser), nil
}

func (r *AuthPostgres) GetUserByCalendarToken(token string) (dto.User, error) {
	var modelUser models.User
	if err := r.db.Where("calendar_token = ?", token).First(&modelUser).Error; err != nil {
		return dto.User{}, err
	}
	return toDTOUser(modelUser), nil
}

func (r *AuthPostgres) SetCalendarToken(userID uuid.UUID, token string) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("calendar_token", token).Error
}

func toDTOUser(m models.User) dto.User {
	return dto.User{
		ID:            m.ID,
		Login:         m.Login,
		Email:         m.Email,
		Name:          m.Name,
		Surname:       m.Surname,
		PasswordHash:  m.PasswordHash,
		FullName:      m.FullName,
		Role:          m.Role,
		CalendarToken: m.CalendarToken,
	}
}
