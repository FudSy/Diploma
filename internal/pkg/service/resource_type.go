package service

import (
	"github.com/FudSy/Diploma/internal/dto"
	"github.com/FudSy/Diploma/internal/pkg/repository"
	"github.com/google/uuid"
)

type ResourceTypeService struct {
	repo repository.ResourceType
}

func NewResourceTypeService(repo repository.ResourceType) *ResourceTypeService {
	return &ResourceTypeService{repo}
}

func (s *ResourceTypeService) Create(name string, options []dto.ResourceTypeOptionRequest) (uuid.UUID, error) {
	return s.repo.Create(name, options)
}

func (s *ResourceTypeService) GetAll() ([]dto.ResourceTypeResponse, error) {
	return s.repo.GetAll()
}

func (s *ResourceTypeService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *ResourceTypeService) AddOption(resourceTypeID uuid.UUID, option dto.ResourceTypeOptionRequest) (uuid.UUID, error) {
	return s.repo.AddOption(resourceTypeID, option)
}

func (s *ResourceTypeService) DeleteOption(optionID uuid.UUID) error {
	return s.repo.DeleteOption(optionID)
}
