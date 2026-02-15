package postgres

import (
	"github.com/FudSy/Diploma/internal/dto"
	"github.com/FudSy/Diploma/internal/pkg/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type BookingPostgres struct {
	db *gorm.DB
}

func NewBookingPostgres(db *gorm.DB) *BookingPostgres {
	return &BookingPostgres{db}
}

func (r *BookingPostgres) Create(userID uuid.UUID, input dto.CreateBookingRequest) (uuid.UUID, error) {
	booking := models.Booking{
		UserID:     userID,
		ResourceID: input.ResourceID,
		StartTime:  input.StartTime,
		EndTime:    input.EndTime,
		Status:     "CONFIRMED",
	}

	if err := r.db.Create(&booking).Error; err != nil {
		return uuid.Nil, err
	}
	return booking.ID, nil
}

func (r *BookingPostgres) GetAllByUser(userID uuid.UUID) ([]dto.BookingResponse, error) {
	var modelBookings []models.Booking
	if err := r.db.Where("user_id = ?", userID).Find(&modelBookings).Error; err != nil {
		return nil, err
	}

	bookings := make([]dto.BookingResponse, 0, len(modelBookings))
	for _, booking := range modelBookings {
		bookings = append(bookings, dto.BookingToResponse(booking))
	}
	return bookings, nil
}

func (r *BookingPostgres) GetById(id uuid.UUID) (dto.BookingResponse, error) {
	var booking models.Booking
	err := r.db.First(&booking, "id = ?", id).Error
	if err != nil {
		return dto.BookingResponse{}, err
	}
	return dto.BookingToResponse(booking), nil
}

func (r *BookingPostgres) HasTimeOverlap(resourceID uuid.UUID, startTime, endTime time.Time) (bool, error) {
	var count int64
	err := r.db.Model(&models.Booking{}).
		Where("resource_id = ? AND status <> ? AND start_time < ? AND end_time > ?", resourceID, "CANCELLED", endTime, startTime).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *BookingPostgres) Update(id uuid.UUID, input dto.UpdateBookingRequest) error {
	updates := map[string]interface{}{}
	if input.StartTime != nil {
		updates["start_time"] = *input.StartTime
	}
	if input.EndTime != nil {
		updates["end_time"] = *input.EndTime
	}
	if input.Status != nil {
		updates["status"] = *input.Status
	}

	if len(updates) == 0 {
		return nil
	}

	return r.db.Model(&models.Booking{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *BookingPostgres) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Booking{}, "id = ?", id).Error
}
