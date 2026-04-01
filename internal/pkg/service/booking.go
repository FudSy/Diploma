package service

import (
	"errors"
	"github.com/FudSy/Diploma/internal/dto"
	"github.com/FudSy/Diploma/internal/pkg/repository"
	"github.com/google/uuid"
	"time"
)

type BookingService struct {
	bookingRepo  repository.Booking
	resourceRepo repository.Resource
}

func NewBookingService(b repository.Booking, r repository.Resource) *BookingService {
	return &BookingService{bookingRepo: b, resourceRepo: r}
}

func (s *BookingService) Create(userID uuid.UUID, input dto.CreateBookingRequest) (uuid.UUID, error) {
	if !input.StartTime.After(time.Now()) {
		return uuid.Nil, errors.New("время начала бронирования должно быть в будущем")
	}

	hasOverlap, err := s.bookingRepo.HasTimeOverlap(input.ResourceID, input.StartTime, input.EndTime)
	if err != nil {
		return uuid.Nil, err
	}
	if hasOverlap {
		return uuid.Nil, errors.New("время бронирования пересекается с существующим")
	}

	return s.bookingRepo.Create(userID, input)
}

func (s *BookingService) GetAll(userID uuid.UUID) ([]dto.BookingResponse, error) {
	return s.bookingRepo.GetAllByUser(userID)
}

func (s *BookingService) GetById(bookingID uuid.UUID) (dto.BookingResponse, error) {
	return s.bookingRepo.GetById(bookingID)
}

func (s *BookingService) Update(userID uuid.UUID, bookingID uuid.UUID, input dto.UpdateBookingRequest) error {
	return s.bookingRepo.Update(bookingID, input)
}

func (s *BookingService) Delete(userID uuid.UUID, bookingID uuid.UUID) error {
	return s.bookingRepo.Delete(bookingID)
}
