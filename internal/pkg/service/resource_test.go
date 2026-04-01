package service

import (
	"errors"
	"testing"

	"github.com/FudSy/Diploma/internal/dto"
	"github.com/google/uuid"
)

type resourceRepoMock struct {
	createFn           func(resource dto.CreateResourceRequest) (uuid.UUID, error)
	getAllFn           func() ([]dto.ResourceResponse, error)
	getByIDFn          func(id uuid.UUID) (dto.ResourceResponse, error)
	updateFn           func(id uuid.UUID, input dto.UpdateResourceRequest) error
	updatePhotoFn      func(id uuid.UUID, photoURL string) error
	increaseCapacityFn func(id uuid.UUID, delta int) error
	decreaseCapacityFn func(id uuid.UUID, delta int) error
	deleteFn           func(id uuid.UUID) error
}

func (m *resourceRepoMock) Create(resource dto.CreateResourceRequest) (uuid.UUID, error) {
	return m.createFn(resource)
}

func (m *resourceRepoMock) GetAll() ([]dto.ResourceResponse, error) {
	return m.getAllFn()
}

func (m *resourceRepoMock) GetById(id uuid.UUID) (dto.ResourceResponse, error) {
	return m.getByIDFn(id)
}

func (m *resourceRepoMock) Update(id uuid.UUID, input dto.UpdateResourceRequest) error {
	return m.updateFn(id, input)
}

func (m *resourceRepoMock) IncreaseCapacity(id uuid.UUID, delta int) error {
	return m.increaseCapacityFn(id, delta)
}

func (m *resourceRepoMock) DecreaseCapacity(id uuid.UUID, delta int) error {
	return m.decreaseCapacityFn(id, delta)
}

func (m *resourceRepoMock) UpdatePhoto(id uuid.UUID, photoURL string) error {
	if m.updatePhotoFn != nil {
		return m.updatePhotoFn(id, photoURL)
	}
	return nil
}

func (m *resourceRepoMock) Delete(id uuid.UUID) error {
	return m.deleteFn(id)
}

func TestResourceService_IncreaseCapacity_InvalidDelta(t *testing.T) {
	svc := NewResourceService(&resourceRepoMock{
		createFn:           func(resource dto.CreateResourceRequest) (uuid.UUID, error) { return uuid.New(), nil },
		getAllFn:           func() ([]dto.ResourceResponse, error) { return nil, nil },
		getByIDFn:          func(id uuid.UUID) (dto.ResourceResponse, error) { return dto.ResourceResponse{}, nil },
		updateFn:           func(id uuid.UUID, input dto.UpdateResourceRequest) error { return nil },
		increaseCapacityFn: func(id uuid.UUID, delta int) error { return nil },
		decreaseCapacityFn: func(id uuid.UUID, delta int) error { return nil },
		deleteFn:           func(id uuid.UUID) error { return nil },
	})

	if err := svc.IncreaseCapacity(uuid.New(), 0); err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestResourceService_DecreaseCapacity_InvalidDelta(t *testing.T) {
	svc := NewResourceService(&resourceRepoMock{
		createFn:           func(resource dto.CreateResourceRequest) (uuid.UUID, error) { return uuid.New(), nil },
		getAllFn:           func() ([]dto.ResourceResponse, error) { return nil, nil },
		getByIDFn:          func(id uuid.UUID) (dto.ResourceResponse, error) { return dto.ResourceResponse{}, nil },
		updateFn:           func(id uuid.UUID, input dto.UpdateResourceRequest) error { return nil },
		increaseCapacityFn: func(id uuid.UUID, delta int) error { return nil },
		decreaseCapacityFn: func(id uuid.UUID, delta int) error { return nil },
		deleteFn:           func(id uuid.UUID) error { return nil },
	})

	if err := svc.DecreaseCapacity(uuid.New(), -1); err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestResourceService_PassThroughCreate(t *testing.T) {
	expectedID := uuid.New()
	svc := NewResourceService(&resourceRepoMock{
		createFn:           func(resource dto.CreateResourceRequest) (uuid.UUID, error) { return expectedID, nil },
		getAllFn:           func() ([]dto.ResourceResponse, error) { return nil, nil },
		getByIDFn:          func(id uuid.UUID) (dto.ResourceResponse, error) { return dto.ResourceResponse{}, nil },
		updateFn:           func(id uuid.UUID, input dto.UpdateResourceRequest) error { return nil },
		increaseCapacityFn: func(id uuid.UUID, delta int) error { return nil },
		decreaseCapacityFn: func(id uuid.UUID, delta int) error { return nil },
		deleteFn:           func(id uuid.UUID) error { return nil },
	})

	id, err := svc.Create(dto.CreateResourceRequest{Name: "Room 1", Type: "MEETING_ROOM", Capacity: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != expectedID {
		t.Fatalf("expected %s, got %s", expectedID, id)
	}
}

func TestResourceService_PassThroughError(t *testing.T) {
	expectedErr := errors.New("repo failed")
	svc := NewResourceService(&resourceRepoMock{
		createFn:           func(resource dto.CreateResourceRequest) (uuid.UUID, error) { return uuid.New(), nil },
		getAllFn:           func() ([]dto.ResourceResponse, error) { return nil, nil },
		getByIDFn:          func(id uuid.UUID) (dto.ResourceResponse, error) { return dto.ResourceResponse{}, nil },
		updateFn:           func(id uuid.UUID, input dto.UpdateResourceRequest) error { return nil },
		increaseCapacityFn: func(id uuid.UUID, delta int) error { return nil },
		decreaseCapacityFn: func(id uuid.UUID, delta int) error { return expectedErr },
		deleteFn:           func(id uuid.UUID) error { return nil },
	})

	err := svc.DecreaseCapacity(uuid.New(), 1)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}
