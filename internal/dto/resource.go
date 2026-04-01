package dto

import "github.com/google/uuid"


type CreateResourceRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Type        string `json:"type" binding:"required"`
	Capacity    int    `json:"capacity" binding:"required,min=1"`
	IsActive    *bool  `json:"is_active"`
	Location    string `json:"location"`
}

type UpdateResourceRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Type        *string `json:"type"`
	Capacity    *int    `json:"capacity"`
	IsActive    *bool   `json:"is_active"`
	Location    *string `json:"location"`
}

type ResourceResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Capacity    int       `json:"capacity"`
	IsActive    bool      `json:"is_active"`
	Location    string    `json:"location"`
	PhotoURL    string    `json:"photo_url"`
}
