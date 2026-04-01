package repository

import (
	"github.com/FudSy/Diploma/internal/dto"
	. "github.com/FudSy/Diploma/internal/pkg/repository/postgres"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Authorization interface {
	CreateUser(user dto.User) (uuid.UUID, error)
	GetUserByLogin(login string) (dto.User, error)
	GetUserById(id uuid.UUID) (dto.User, error)
}

type Resource interface {
	Create(resource dto.CreateResourceRequest) (uuid.UUID, error)
	GetAll() ([]dto.ResourceResponse, error)
	GetById(id uuid.UUID) (dto.ResourceResponse, error)
	Update(id uuid.UUID, input dto.UpdateResourceRequest) error
	UpdatePhoto(id uuid.UUID, photoURL string) error
	IncreaseCapacity(id uuid.UUID, delta int) error
	DecreaseCapacity(id uuid.UUID, delta int) error
	Delete(id uuid.UUID) error
}

type Booking interface {
	Create(userID uuid.UUID, input dto.CreateBookingRequest) (uuid.UUID, error)
	GetAllByUser(userID uuid.UUID) ([]dto.BookingResponse, error)
	GetById(id uuid.UUID) (dto.BookingResponse, error)
	HasTimeOverlap(resourceID uuid.UUID, startTime, endTime time.Time) (bool, error)
	Update(id uuid.UUID, input dto.UpdateBookingRequest) error
	Delete(id uuid.UUID) error
}

type ResourceType interface {
	Create(name string, options []dto.ResourceTypeOptionRequest) (uuid.UUID, error)
	GetAll() ([]dto.ResourceTypeResponse, error)
	Delete(id uuid.UUID) error
	AddOption(resourceTypeID uuid.UUID, option dto.ResourceTypeOptionRequest) (uuid.UUID, error)
	DeleteOption(optionID uuid.UUID) error
}

type Repository struct {
	Authorization
	Resource
	Booking
	ResourceType
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Resource:      NewResourcePostgres(db),
		Booking:       NewBookingPostgres(db),
		ResourceType:  NewResourceTypePostgres(db),
	}
}
