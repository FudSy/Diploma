package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/FudSy/Diploma/internal/dto"
	"github.com/FudSy/Diploma/internal/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type authSvcMock struct {
	createUserFn   func(input dto.RegisterRequest) (uuid.UUID, error)
	createAdminFn  func(input dto.RegisterRequest) (uuid.UUID, error)
	loginFn        func(login, password string) (string, error)
	parseTokenFn   func(accessToken string) (uuid.UUID, error)
	getUserByIDFn  func(id uuid.UUID) (dto.User, error)
	checkPasswordf func(user dto.User, password string) error
}

func (m *authSvcMock) CreateUser(input dto.RegisterRequest) (uuid.UUID, error) {
	return m.createUserFn(input)
}
func (m *authSvcMock) CreateAdmin(input dto.RegisterRequest) (uuid.UUID, error) {
	return m.createAdminFn(input)
}
func (m *authSvcMock) Login(login, password string) (string, error) {
	return m.loginFn(login, password)
}
func (m *authSvcMock) ParseToken(accessToken string) (uuid.UUID, error) {
	return m.parseTokenFn(accessToken)
}
func (m *authSvcMock) GetUserById(id uuid.UUID) (dto.User, error) {
	return m.getUserByIDFn(id)
}
func (m *authSvcMock) CheckPassword(user dto.User, password string) error {
	return m.checkPasswordf(user, password)
}

type resourceSvcMock struct {
	createFn           func(input dto.CreateResourceRequest) (uuid.UUID, error)
	getAllFn           func() ([]dto.ResourceResponse, error)
	getByIDFn          func(id uuid.UUID) (dto.ResourceResponse, error)
	deleteFn           func(id uuid.UUID) error
	updateFn           func(id uuid.UUID, input dto.UpdateResourceRequest) error
	increaseCapacityFn func(id uuid.UUID, delta int) error
	decreaseCapacityFn func(id uuid.UUID, delta int) error
}

func (m *resourceSvcMock) Create(input dto.CreateResourceRequest) (uuid.UUID, error) {
	return m.createFn(input)
}
func (m *resourceSvcMock) GetAll() ([]dto.ResourceResponse, error) {
	return m.getAllFn()
}
func (m *resourceSvcMock) GetById(id uuid.UUID) (dto.ResourceResponse, error) {
	return m.getByIDFn(id)
}
func (m *resourceSvcMock) Delete(id uuid.UUID) error {
	return m.deleteFn(id)
}
func (m *resourceSvcMock) Update(id uuid.UUID, input dto.UpdateResourceRequest) error {
	return m.updateFn(id, input)
}
func (m *resourceSvcMock) IncreaseCapacity(id uuid.UUID, delta int) error {
	return m.increaseCapacityFn(id, delta)
}
func (m *resourceSvcMock) DecreaseCapacity(id uuid.UUID, delta int) error {
	return m.decreaseCapacityFn(id, delta)
}

type bookingSvcMock struct {
	createFn  func(userID uuid.UUID, input dto.CreateBookingRequest) (uuid.UUID, error)
	getAllFn  func(userID uuid.UUID) ([]dto.BookingResponse, error)
	getByIDFn func(bookingID uuid.UUID) (dto.BookingResponse, error)
	updateFn  func(userID, bookingID uuid.UUID, input dto.UpdateBookingRequest) error
	deleteFn  func(userID, bookingID uuid.UUID) error
}

func (m *bookingSvcMock) Create(userID uuid.UUID, input dto.CreateBookingRequest) (uuid.UUID, error) {
	return m.createFn(userID, input)
}
func (m *bookingSvcMock) GetAll(userID uuid.UUID) ([]dto.BookingResponse, error) {
	return m.getAllFn(userID)
}
func (m *bookingSvcMock) GetById(bookingID uuid.UUID) (dto.BookingResponse, error) {
	return m.getByIDFn(bookingID)
}
func (m *bookingSvcMock) Update(userID, bookingID uuid.UUID, input dto.UpdateBookingRequest) error {
	return m.updateFn(userID, bookingID, input)
}
func (m *bookingSvcMock) Delete(userID, bookingID uuid.UUID) error {
	return m.deleteFn(userID, bookingID)
}

func newTestHandler(userID, adminID uuid.UUID, bookingOwnerID uuid.UUID) *Handler {
	auth := &authSvcMock{
		createUserFn:  func(input dto.RegisterRequest) (uuid.UUID, error) { return uuid.New(), nil },
		createAdminFn: func(input dto.RegisterRequest) (uuid.UUID, error) { return uuid.New(), nil },
		loginFn:       func(login, password string) (string, error) { return "token-user", nil },
		parseTokenFn: func(accessToken string) (uuid.UUID, error) {
			switch accessToken {
			case "token-user":
				return userID, nil
			case "token-admin":
				return adminID, nil
			default:
				return uuid.Nil, errors.New("invalid token")
			}
		},
		getUserByIDFn: func(id uuid.UUID) (dto.User, error) {
			switch id {
			case userID:
				return dto.User{
					ID:      userID,
					Role:    "USER",
					Email:   "user@example.com",
					Name:    "Ivan",
					Surname: "Petrov",
				}, nil
			case adminID:
				return dto.User{
					ID:      adminID,
					Role:    "ADMIN",
					Email:   "admin@example.com",
					Name:    "Anna",
					Surname: "Smirnova",
				}, nil
			default:
				return dto.User{}, errors.New("user not found")
			}
		},
		checkPasswordf: func(user dto.User, password string) error { return nil },
	}
	res := &resourceSvcMock{
		createFn:           func(input dto.CreateResourceRequest) (uuid.UUID, error) { return uuid.New(), nil },
		getAllFn:           func() ([]dto.ResourceResponse, error) { return []dto.ResourceResponse{}, nil },
		getByIDFn:          func(id uuid.UUID) (dto.ResourceResponse, error) { return dto.ResourceResponse{ID: id}, nil },
		deleteFn:           func(id uuid.UUID) error { return nil },
		updateFn:           func(id uuid.UUID, input dto.UpdateResourceRequest) error { return nil },
		increaseCapacityFn: func(id uuid.UUID, delta int) error { return nil },
		decreaseCapacityFn: func(id uuid.UUID, delta int) error { return nil },
	}
	book := &bookingSvcMock{
		createFn: func(userID uuid.UUID, input dto.CreateBookingRequest) (uuid.UUID, error) {
			return uuid.New(), nil
		},
		getAllFn: func(userID uuid.UUID) ([]dto.BookingResponse, error) { return []dto.BookingResponse{}, nil },
		getByIDFn: func(bookingID uuid.UUID) (dto.BookingResponse, error) {
			if bookingID == uuid.Nil {
				return dto.BookingResponse{}, gorm.ErrRecordNotFound
			}
			return dto.BookingResponse{
				ID:         bookingID,
				UserID:     bookingOwnerID,
				ResourceID: uuid.New(),
				StartTime:  time.Now().Add(time.Hour),
				EndTime:    time.Now().Add(2 * time.Hour),
				Status:     "CONFIRMED",
			}, nil
		},
		updateFn: func(userID, bookingID uuid.UUID, input dto.UpdateBookingRequest) error { return nil },
		deleteFn: func(userID, bookingID uuid.UUID) error { return nil },
	}

	return NewHandler(&service.Service{
		Authorization: auth,
		Resource:      res,
		Booking:       book,
	})
}

func performRequest(r http.Handler, method, path string, body []byte, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestInitRoutes_PublicRegisterWithoutToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := newTestHandler(uuid.New(), uuid.New(), uuid.New())

	w := performRequest(h.InitRoutes(), http.MethodPost, "/auth/register", []byte(`{}`), "")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestInitRoutes_ProtectedEndpointRequiresToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := newTestHandler(uuid.New(), uuid.New(), uuid.New())

	w := performRequest(h.InitRoutes(), http.MethodGet, "/auth/me", nil, "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestInitRoutes_MeReturnsNameAndSurname(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	adminID := uuid.New()
	h := newTestHandler(userID, adminID, userID)

	w := performRequest(h.InitRoutes(), http.MethodGet, "/auth/me", nil, "token-user")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp dto.MeResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Name != "Ivan" {
		t.Fatalf("expected name Ivan, got %q", resp.Name)
	}
	if resp.Surname != "Petrov" {
		t.Fatalf("expected surname Petrov, got %q", resp.Surname)
	}
}

func TestInitRoutes_ResourcesAccessByRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	adminID := uuid.New()
	h := newTestHandler(userID, adminID, userID)
	r := h.InitRoutes()

	getW := performRequest(r, http.MethodGet, "/resources/", nil, "token-user")
	if getW.Code != http.StatusOK {
		t.Fatalf("expected user GET /resources to be 200, got %d", getW.Code)
	}

	createBody, _ := json.Marshal(dto.CreateResourceRequest{
		Name: "Room A", Type: "MEETING_ROOM", Capacity: 3,
	})
	userPost := performRequest(r, http.MethodPost, "/resources/", createBody, "token-user")
	if userPost.Code != http.StatusForbidden {
		t.Fatalf("expected user POST /resources to be 403, got %d", userPost.Code)
	}

	adminPost := performRequest(r, http.MethodPost, "/resources/", createBody, "token-admin")
	if adminPost.Code != http.StatusOK {
		t.Fatalf("expected admin POST /resources to be 200, got %d", adminPost.Code)
	}
}

func TestInitRoutes_AdminEndpointDeniedForUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	adminID := uuid.New()
	h := newTestHandler(userID, adminID, userID)

	w := performRequest(h.InitRoutes(), http.MethodGet, "/auth/admin/check", nil, "token-user")
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestInitRoutes_BookingOwnership(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	adminID := uuid.New()
	ownerID := uuid.New()
	bookingID := uuid.New()
	h := newTestHandler(userID, adminID, ownerID)
	r := h.InitRoutes()

	userW := performRequest(r, http.MethodGet, "/bookings/"+bookingID.String(), nil, "token-user")
	if userW.Code != http.StatusForbidden {
		t.Fatalf("expected non-owner USER to get 403, got %d", userW.Code)
	}

	adminW := performRequest(r, http.MethodGet, "/bookings/"+bookingID.String(), nil, "token-admin")
	if adminW.Code != http.StatusOK {
		t.Fatalf("expected ADMIN to get 200, got %d", adminW.Code)
	}
}
