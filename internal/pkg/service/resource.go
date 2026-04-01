package service

import (
	"errors"

	"github.com/FudSy/Diploma/internal/dto"
	"github.com/FudSy/Diploma/internal/pkg/repository"
	"github.com/google/uuid"
)

type ResourceService struct {
	repo repository.Resource
}

func NewResourceService(repo repository.Resource) *ResourceService {
	return &ResourceService{repo: repo}
}

func (s *ResourceService) Create(input dto.CreateResourceRequest) (uuid.UUID, error) {
	return s.repo.Create(input)
}

func (s *ResourceService) GetAll() ([]dto.ResourceResponse, error) {
	return s.repo.GetAll()
}

func (s *ResourceService) GetById(id uuid.UUID) (dto.ResourceResponse, error) {
	return s.repo.GetById(id)
}

func (s *ResourceService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *ResourceService) Update(id uuid.UUID, input dto.UpdateResourceRequest) error {
	return s.repo.Update(id, input)
}

func (s *ResourceService) UpdatePhoto(id uuid.UUID, photoURL string) error {
	return s.repo.UpdatePhoto(id, photoURL)
}

func (s *ResourceService) IncreaseCapacity(id uuid.UUID, delta int) error {
	if delta <= 0 {
		return errors.New("прирост вместимости должен быть больше нуля")
	}

	return s.repo.IncreaseCapacity(id, delta)
}

func (s *ResourceService) DecreaseCapacity(id uuid.UUID, delta int) error {
	if delta <= 0 {
		return errors.New("уменьшение вместимости должно быть больше нуля")
	}

	return s.repo.DecreaseCapacity(id, delta)
}
