package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FudSy/Diploma/internal/dto"
	"github.com/FudSy/Diploma/internal/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestUserIdentity_RejectsInvalidBearerHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	adminID := uuid.New()
	h := newTestHandler(userID, adminID, userID)

	r := gin.New()
	r.GET("/protected", h.userIdentity, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Token abc")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAdminIdentity_RequiresAdminRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	adminID := uuid.New()
	h := newTestHandler(userID, adminID, userID)

	r := gin.New()
	r.GET("/admin", h.userIdentity, h.adminIdentity, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	userReq := httptest.NewRequest(http.MethodGet, "/admin", nil)
	userReq.Header.Set("Authorization", "Bearer token-user")
	userW := httptest.NewRecorder()
	r.ServeHTTP(userW, userReq)
	if userW.Code != http.StatusForbidden {
		t.Fatalf("expected user 403, got %d", userW.Code)
	}

	adminReq := httptest.NewRequest(http.MethodGet, "/admin", nil)
	adminReq.Header.Set("Authorization", "Bearer token-admin")
	adminW := httptest.NewRecorder()
	r.ServeHTTP(adminW, adminReq)
	if adminW.Code != http.StatusOK {
		t.Fatalf("expected admin 200, got %d", adminW.Code)
	}
}

func TestGetUserID_ContextErrors(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	if _, err := getUserID(c); err == nil {
		t.Fatal("expected error when user id is missing")
	}

	c.Set(userCtx, "not-uuid")
	if _, err := getUserID(c); err == nil {
		t.Fatal("expected error when user id has invalid type")
	}
}

func TestUserIdentity_ParseTokenError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	auth := &authSvcMock{
		createUserFn:   func(input dto.RegisterRequest) (uuid.UUID, error) { return uuid.New(), nil },
		createAdminFn:  func(input dto.RegisterRequest) (uuid.UUID, error) { return uuid.New(), nil },
		loginFn:        func(login, password string) (string, error) { return "token", nil },
		parseTokenFn:   func(accessToken string) (uuid.UUID, error) { return uuid.Nil, errors.New("bad token") },
		getUserByIDFn:  func(id uuid.UUID) (dto.User, error) { return dto.User{}, nil },
		checkPasswordf: func(user dto.User, password string) error { return nil },
	}
	h := NewHandler(&service.Service{
		Authorization: auth,
		Resource: &resourceSvcMock{
			createFn:           func(input dto.CreateResourceRequest) (uuid.UUID, error) { return uuid.New(), nil },
			getAllFn:           func() ([]dto.ResourceResponse, error) { return nil, nil },
			getByIDFn:          func(id uuid.UUID) (dto.ResourceResponse, error) { return dto.ResourceResponse{}, nil },
			deleteFn:           func(id uuid.UUID) error { return nil },
			updateFn:           func(id uuid.UUID, input dto.UpdateResourceRequest) error { return nil },
			increaseCapacityFn: func(id uuid.UUID, delta int) error { return nil },
			decreaseCapacityFn: func(id uuid.UUID, delta int) error { return nil },
		},
		Booking: &bookingSvcMock{
			createFn:  func(userID uuid.UUID, input dto.CreateBookingRequest) (uuid.UUID, error) { return uuid.New(), nil },
			getAllFn:  func(userID uuid.UUID) ([]dto.BookingResponse, error) { return nil, nil },
			getByIDFn: func(bookingID uuid.UUID) (dto.BookingResponse, error) { return dto.BookingResponse{}, nil },
			updateFn:  func(userID, bookingID uuid.UUID, input dto.UpdateBookingRequest) error { return nil },
			deleteFn:  func(userID, bookingID uuid.UUID) error { return nil },
		},
	})

	r := gin.New()
	r.GET("/protected", h.userIdentity, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestUserIdentity_UsesCookieWhenNoHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	adminID := uuid.New()
	h := newTestHandler(userID, adminID, userID)

	r := gin.New()
	r.GET("/protected", h.userIdentity, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{Name: authCookieName, Value: "token-user"})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
