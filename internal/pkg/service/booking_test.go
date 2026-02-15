package service

import (
	"errors"
	"testing"
	"time"

	"github.com/FudSy/Diploma/internal/dto"
	"github.com/google/uuid"
)

type bookingRepoMock struct {
	createFn         func(userID uuid.UUID, input dto.CreateBookingRequest) (uuid.UUID, error)
	getAllByUserFn   func(userID uuid.UUID) ([]dto.BookingResponse, error)
	getByIDFn        func(id uuid.UUID) (dto.BookingResponse, error)
	hasTimeOverlapFn func(resourceID uuid.UUID, startTime, endTime time.Time) (bool, error)
	updateFn         func(id uuid.UUID, input dto.UpdateBookingRequest) error
	deleteFn         func(id uuid.UUID) error
}

func (m *bookingRepoMock) Create(userID uuid.UUID, input dto.CreateBookingRequest) (uuid.UUID, error) {
	return m.createFn(userID, input)
}

func (m *bookingRepoMock) GetAllByUser(userID uuid.UUID) ([]dto.BookingResponse, error) {
	return m.getAllByUserFn(userID)
}

func (m *bookingRepoMock) GetById(id uuid.UUID) (dto.BookingResponse, error) {
	return m.getByIDFn(id)
}

func (m *bookingRepoMock) HasTimeOverlap(resourceID uuid.UUID, startTime, endTime time.Time) (bool, error) {
	return m.hasTimeOverlapFn(resourceID, startTime, endTime)
}

func (m *bookingRepoMock) Update(id uuid.UUID, input dto.UpdateBookingRequest) error {
	return m.updateFn(id, input)
}

func (m *bookingRepoMock) Delete(id uuid.UUID) error {
	return m.deleteFn(id)
}

func TestBookingService_Create_StartTimeMustBeFuture(t *testing.T) {
	svc := NewBookingService(
		&bookingRepoMock{
			createFn:         func(userID uuid.UUID, input dto.CreateBookingRequest) (uuid.UUID, error) { return uuid.New(), nil },
			getAllByUserFn:   func(userID uuid.UUID) ([]dto.BookingResponse, error) { return nil, nil },
			getByIDFn:        func(id uuid.UUID) (dto.BookingResponse, error) { return dto.BookingResponse{}, nil },
			hasTimeOverlapFn: func(resourceID uuid.UUID, startTime, endTime time.Time) (bool, error) { return false, nil },
			updateFn:         func(id uuid.UUID, input dto.UpdateBookingRequest) error { return nil },
			deleteFn:         func(id uuid.UUID) error { return nil },
		},
		nil,
	)

	_, err := svc.Create(uuid.New(), dto.CreateBookingRequest{
		ResourceID: uuid.New(),
		StartTime:  time.Now().Add(-time.Hour),
		EndTime:    time.Now().Add(time.Hour),
	})
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestBookingService_Create_Overlap(t *testing.T) {
	svc := NewBookingService(
		&bookingRepoMock{
			createFn:       func(userID uuid.UUID, input dto.CreateBookingRequest) (uuid.UUID, error) { return uuid.New(), nil },
			getAllByUserFn: func(userID uuid.UUID) ([]dto.BookingResponse, error) { return nil, nil },
			getByIDFn:      func(id uuid.UUID) (dto.BookingResponse, error) { return dto.BookingResponse{}, nil },
			hasTimeOverlapFn: func(resourceID uuid.UUID, startTime, endTime time.Time) (bool, error) {
				return true, nil
			},
			updateFn: func(id uuid.UUID, input dto.UpdateBookingRequest) error { return nil },
			deleteFn: func(id uuid.UUID) error { return nil },
		},
		nil,
	)

	_, err := svc.Create(uuid.New(), dto.CreateBookingRequest{
		ResourceID: uuid.New(),
		StartTime:  time.Now().Add(time.Hour),
		EndTime:    time.Now().Add(2 * time.Hour),
	})
	if err == nil {
		t.Fatal("expected overlap error, got nil")
	}
}

func TestBookingService_Create_Success(t *testing.T) {
	expectedID := uuid.New()
	svc := NewBookingService(
		&bookingRepoMock{
			createFn: func(userID uuid.UUID, input dto.CreateBookingRequest) (uuid.UUID, error) {
				return expectedID, nil
			},
			getAllByUserFn:   func(userID uuid.UUID) ([]dto.BookingResponse, error) { return nil, nil },
			getByIDFn:        func(id uuid.UUID) (dto.BookingResponse, error) { return dto.BookingResponse{}, nil },
			hasTimeOverlapFn: func(resourceID uuid.UUID, startTime, endTime time.Time) (bool, error) { return false, nil },
			updateFn:         func(id uuid.UUID, input dto.UpdateBookingRequest) error { return nil },
			deleteFn:         func(id uuid.UUID) error { return nil },
		},
		nil,
	)

	id, err := svc.Create(uuid.New(), dto.CreateBookingRequest{
		ResourceID: uuid.New(),
		StartTime:  time.Now().Add(time.Hour),
		EndTime:    time.Now().Add(2 * time.Hour),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != expectedID {
		t.Fatalf("expected %s, got %s", expectedID, id)
	}
}

func TestBookingService_PassThroughMethods(t *testing.T) {
	userID := uuid.New()
	bookingID := uuid.New()
	repoErr := errors.New("repo error")
	svc := NewBookingService(
		&bookingRepoMock{
			createFn:         func(userID uuid.UUID, input dto.CreateBookingRequest) (uuid.UUID, error) { return uuid.New(), nil },
			getAllByUserFn:   func(gotUserID uuid.UUID) ([]dto.BookingResponse, error) { return nil, repoErr },
			getByIDFn:        func(id uuid.UUID) (dto.BookingResponse, error) { return dto.BookingResponse{}, repoErr },
			hasTimeOverlapFn: func(resourceID uuid.UUID, startTime, endTime time.Time) (bool, error) { return false, nil },
			updateFn:         func(id uuid.UUID, input dto.UpdateBookingRequest) error { return repoErr },
			deleteFn:         func(id uuid.UUID) error { return repoErr },
		},
		nil,
	)

	if _, err := svc.GetAll(userID); !errors.Is(err, repoErr) {
		t.Fatalf("expected %v, got %v", repoErr, err)
	}
	if _, err := svc.GetById(bookingID); !errors.Is(err, repoErr) {
		t.Fatalf("expected %v, got %v", repoErr, err)
	}
	if err := svc.Update(userID, bookingID, dto.UpdateBookingRequest{}); !errors.Is(err, repoErr) {
		t.Fatalf("expected %v, got %v", repoErr, err)
	}
	if err := svc.Delete(userID, bookingID); !errors.Is(err, repoErr) {
		t.Fatalf("expected %v, got %v", repoErr, err)
	}
}
