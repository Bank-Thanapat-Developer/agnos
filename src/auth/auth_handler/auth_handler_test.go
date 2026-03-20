package authhandler

import (
	"agnos-test/dto"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockAuthUsecase struct {
	loginFunc func(req dto.LoginRequest, hospitalID uint) (*dto.LoginResponse, error)
}

func (m *mockAuthUsecase) Login(req dto.LoginRequest, hospitalID uint) (*dto.LoginResponse, error) {
	return m.loginFunc(req, hospitalID)
}

func setupRouter(handler *AuthHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/login", func(c *gin.Context) {
		c.Set("hospital_id", uint(1))
		c.Next()
	}, handler.Login)
	return r
}

func TestHandler_Login_Success(t *testing.T) {
	uc := &mockAuthUsecase{
		loginFunc: func(req dto.LoginRequest, hospitalID uint) (*dto.LoginResponse, error) {
			return &dto.LoginResponse{
				Token:        "test-token",
				StaffID:      1,
				Username:     req.Username,
				HospitalID:   hospitalID,
				HospitalName: "Hospital A",
			}, nil
		},
	}
	handler := NewAuthHandler(uc)
	r := setupRouter(handler)

	body, _ := json.Marshal(dto.LoginRequest{Username: "admin", Password: "password123"})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.LoginResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "test-token", resp.Token)
	assert.Equal(t, "admin", resp.Username)
}

func TestHandler_Login_InvalidBody(t *testing.T) {
	uc := &mockAuthUsecase{}
	handler := NewAuthHandler(uc)
	r := setupRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_Login_Unauthorized(t *testing.T) {
	uc := &mockAuthUsecase{
		loginFunc: func(req dto.LoginRequest, hospitalID uint) (*dto.LoginResponse, error) {
			return nil, errors.New("invalid password")
		},
	}
	handler := NewAuthHandler(uc)
	r := setupRouter(handler)

	body, _ := json.Marshal(dto.LoginRequest{Username: "admin", Password: "wrong"})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
