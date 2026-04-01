package service

import (
	"github.com/FudSy/Diploma/internal/dto"
	"github.com/FudSy/Diploma/internal/pkg/repository"
	"github.com/google/uuid"
)

type Authorization interface {
	CreateUser(input dto.RegisterRequest) (uuid.UUID, error)
	CreateAdmin(input dto.RegisterRequest) (uuid.UUID, error)
	Login(login, password string) (string, error)
	ParseToken(accessToken string) (uuid.UUID, error)
	GetUserById(id uuid.UUID) (dto.User, error)
	CheckPassword(user dto.User, password string) error
}

type Resource interface {
	Create(input dto.CreateResourceRequest) (uuid.UUID, error)
	GetAll() ([]dto.ResourceResponse, error)
	GetById(id uuid.UUID) (dto.ResourceResponse, error)
	Delete(id uuid.UUID) error
	Update(id uuid.UUID, input dto.UpdateResourceRequest) error
	UpdatePhoto(id uuid.UUID, photoURL string) error
	IncreaseCapacity(id uuid.UUID, delta int) error
	DecreaseCapacity(id uuid.UUID, delta int) error
}

type Booking interface {
	Create(userID uuid.UUID, input dto.CreateBookingRequest) (uuid.UUID, error)
	GetAll(userID uuid.UUID) ([]dto.BookingResponse, error)
	GetById(bookingID uuid.UUID) (dto.BookingResponse, error)
	Update(userID, bookingID uuid.UUID, input dto.UpdateBookingRequest) error
	Delete(userID, bookingID uuid.UUID) error
}

type ResourceType interface {
	Create(name string, options []dto.ResourceTypeOptionRequest) (uuid.UUID, error)
	GetAll() ([]dto.ResourceTypeResponse, error)
	Delete(id uuid.UUID) error
	AddOption(resourceTypeID uuid.UUID, option dto.ResourceTypeOptionRequest) (uuid.UUID, error)
	DeleteOption(optionID uuid.UUID) error
}

type Service struct {
	Authorization
	Resource
	Booking
	ResourceType
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Resource:      NewResourceService(repos.Resource),
		Booking:       NewBookingService(repos.Booking, repos.Resource),
		ResourceType:  NewResourceTypeService(repos.ResourceType),
	}
}
