package postgres

import (
	"errors"

	"github.com/FudSy/Diploma/internal/dto"
	"github.com/FudSy/Diploma/internal/pkg/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ResourcePostgres struct {
	db *gorm.DB
}

func NewResourcePostgres(db *gorm.DB) *ResourcePostgres {
	return &ResourcePostgres{db}
}

func (r *ResourcePostgres) Create(resource dto.CreateResourceRequest) (uuid.UUID, error) {
	modelResource := models.Resource{
		Name:        resource.Name,
		Description: resource.Description,
		Type:        resource.Type,
		Capacity:    resource.Capacity,
		IsActive:    true,
		Location:    resource.Location,
	}
	if resource.IsActive != nil {
		modelResource.IsActive = *resource.IsActive
	}

	if err := r.db.Create(&modelResource).Error; err != nil {
		return uuid.Nil, err
	}
	return modelResource.ID, nil
}

func (r *ResourcePostgres) GetAll() ([]dto.ResourceResponse, error) {
	var modelResources []models.Resource
	if err := r.db.Find(&modelResources).Error; err != nil {
		return nil, err
	}

	resources := make([]dto.ResourceResponse, 0, len(modelResources))
	for _, resource := range modelResources {
		resources = append(resources, dto.ResourceToResponse(resource))
	}
	return resources, nil
}

func (r *ResourcePostgres) GetById(id uuid.UUID) (dto.ResourceResponse, error) {
	var resource models.Resource
	err := r.db.First(&resource, "id = ?", id).Error
	if err != nil {
		return dto.ResourceResponse{}, err
	}
	return dto.ResourceToResponse(resource), nil
}

func (r *ResourcePostgres) Update(id uuid.UUID, input dto.UpdateResourceRequest) error {
	updates := map[string]interface{}{}
	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Description != nil {
		updates["description"] = *input.Description
	}
	if input.Type != nil {
		updates["type"] = *input.Type
	}
	if input.Capacity != nil {
		updates["capacity"] = *input.Capacity
	}
	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
	}
	if input.Location != nil {
		updates["location"] = *input.Location
	}

	if len(updates) == 0 {
		return nil
	}

	return r.db.Model(&models.Resource{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *ResourcePostgres) IncreaseCapacity(id uuid.UUID, delta int) error {
	result := r.db.Model(&models.Resource{}).
		Where("id = ?", id).
		Update("capacity", gorm.Expr("capacity + ?", delta))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *ResourcePostgres) DecreaseCapacity(id uuid.UUID, delta int) error {
	result := r.db.Model(&models.Resource{}).
		Where("id = ? AND capacity >= ?", id, delta).
		Update("capacity", gorm.Expr("capacity - ?", delta))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected > 0 {
		return nil
	}

	var count int64
	if err := r.db.Model(&models.Resource{}).
		Where("id = ?", id).
		Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return gorm.ErrRecordNotFound
	}

	return errors.New("недостаточная вместимость ресурса")
}

func (r *ResourcePostgres) UpdatePhoto(id uuid.UUID, photoURL string) error {
	result := r.db.Model(&models.Resource{}).
		Where("id = ?", id).
		Update("photo_url", photoURL)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *ResourcePostgres) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Resource{}, "id = ?", id).Error
}
