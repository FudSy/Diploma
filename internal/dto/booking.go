package dto

import (
	"time"
	"github.com/google/uuid"
)

// CreateBookingRequest
type CreateBookingRequest struct {
	ResourceID uuid.UUID `json:"resource_id" binding:"required"`
	StartTime  time.Time `json:"start_time" binding:"required"`
	EndTime    time.Time `json:"end_time" binding:"required,gtfield=StartTime"`
}

// UpdateBookingRequest
type UpdateBookingRequest struct {
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
	Status    *string    `json:"status"` // CONFIRMED, CANCELLED
}

// BookingResponse
type BookingResponse struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	ResourceID uuid.UUID `json:"resource_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Status     string    `json:"status"`
}

// AdminBookingResponse includes user and resource names for admin view
type AdminBookingResponse struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	UserName     string    `json:"user_name"`
	ResourceID   uuid.UUID `json:"resource_id"`
	ResourceName string    `json:"resource_name"`
	ResourceType string    `json:"resource_type"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Status       string    `json:"status"`
}

// BusySlot represents a booked time interval for a resource
type BusySlot struct {
	BookingID uuid.UUID `json:"booking_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    string    `json:"status"`
}
