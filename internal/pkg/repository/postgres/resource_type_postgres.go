package postgres

import (
	"github.com/FudSy/Diploma/internal/dto"
	"github.com/FudSy/Diploma/internal/pkg/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ResourceTypePostgres struct {
	db *gorm.DB
}

func NewResourceTypePostgres(db *gorm.DB) *ResourceTypePostgres {
	return &ResourceTypePostgres{db}
}

func (r *ResourceTypePostgres) Create(name string, options []dto.ResourceTypeOptionRequest) (uuid.UUID, error) {
	rt := models.ResourceType{Name: name}

	for _, o := range options {
		rt.Options = append(rt.Options, models.ResourceTypeOption{
			Name:       o.Name,
			OptionType: o.OptionType,
			IsRequired: o.IsRequired,
		})
	}

	if err := r.db.Create(&rt).Error; err != nil {
		return uuid.Nil, err
	}
	return rt.ID, nil
}

func (r *ResourceTypePostgres) GetAll() ([]dto.ResourceTypeResponse, error) {
	var types []models.ResourceType
	if err := r.db.Preload("Options").Find(&types).Error; err != nil {
		return nil, err
	}
	result := make([]dto.ResourceTypeResponse, 0, len(types))
	for _, t := range types {
		opts := make([]dto.ResourceTypeOptionResponse, 0, len(t.Options))
		for _, o := range t.Options {
			opts = append(opts, dto.ResourceTypeOptionResponse{
				ID:         o.ID,
				Name:       o.Name,
				OptionType: o.OptionType,
				IsRequired: o.IsRequired,
			})
		}
		result = append(result, dto.ResourceTypeResponse{ID: t.ID, Name: t.Name, Options: opts})
	}
	return result, nil
}

func (r *ResourceTypePostgres) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.ResourceType{}, "id = ?", id).Error
}

func (r *ResourceTypePostgres) AddOption(resourceTypeID uuid.UUID, option dto.ResourceTypeOptionRequest) (uuid.UUID, error) {
	o := models.ResourceTypeOption{
		ResourceTypeID: resourceTypeID,
		Name:           option.Name,
		OptionType:     option.OptionType,
		IsRequired:     option.IsRequired,
	}
	if err := r.db.Create(&o).Error; err != nil {
		return uuid.Nil, err
	}
	return o.ID, nil
}

func (r *ResourceTypePostgres) DeleteOption(optionID uuid.UUID) error {
	return r.db.Delete(&models.ResourceTypeOption{}, "id = ?", optionID).Error
}
