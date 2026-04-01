package dto

import "github.com/google/uuid"

type ResourceTypeOptionRequest struct {
	Name       string `json:"name" binding:"required"`
	OptionType string `json:"option_type" binding:"required,oneof=text number boolean"`
	IsRequired bool   `json:"is_required"`
}

type ResourceTypeOptionResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	OptionType string    `json:"option_type"`
	IsRequired bool      `json:"is_required"`
}

type CreateResourceTypeRequest struct {
	Name    string                      `json:"name" binding:"required"`
	Options []ResourceTypeOptionRequest `json:"options"`
}

type ResourceTypeResponse struct {
	ID      uuid.UUID                    `json:"id"`
	Name    string                       `json:"name"`
	Options []ResourceTypeOptionResponse `json:"options"`
}
