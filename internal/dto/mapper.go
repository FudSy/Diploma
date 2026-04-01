package dto

import (
	"github.com/FudSy/Diploma/internal/pkg/models"
)

// BookingToResponse конвертирует модель Booking в DTO
func BookingToResponse(b models.Booking) BookingResponse {
	return BookingResponse{
		ID:         b.ID,
		UserID:     b.UserID,
		ResourceID: b.ResourceID,
		StartTime:  b.StartTime,
		EndTime:    b.EndTime,
		Status:     b.Status,
	}
}

// ResourceToResponse конвертирует модель Resource в DTO
func ResourceToResponse(r models.Resource) ResourceResponse {
	return ResourceResponse{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Type:        r.Type,
		Capacity:    r.Capacity,
		IsActive:    r.IsActive,
		Location:    r.Location,
		PhotoURL:    r.PhotoURL,
	}
}

// UserToMeResponse конвертирует модель User в DTO для /auth/me
func UserToMeResponse(u models.User) MeResponse {
	return MeResponse{
		ID:      u.ID,
		Email:   u.Email,
		Name:    u.Name,
		Surname: u.Surname,
		Role:    u.Role,
	}
}
