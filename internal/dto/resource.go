package dto

import "github.com/google/uuid"


type CreateResourceRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Type        string `json:"type" binding:"required,oneof=MEETING_ROOM CAR DEVICE"`
	Capacity    int    `json:"capacity" binding:"required,min=1"`
	IsActive    *bool  `json:"is_active"`
}

type UpdateResourceRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Type        *string `json:"type"`
	Capacity    *int    `json:"capacity"`
	IsActive    *bool   `json:"is_active"`
}

type ResourceResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Capacity    int       `json:"capacity"`
	IsActive    bool      `json:"is_active"`
}
