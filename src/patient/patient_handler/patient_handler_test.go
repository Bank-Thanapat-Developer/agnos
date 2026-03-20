package patienthandler

import (
	"agnos-test/dto"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockPatientUsecase struct {
	findFunc   func(id string) (dto.PatientSearchResponse, error)
	createFunc func(req dto.CreatePatientRequest, staffID uint, hospitalID uint) (*dto.CreatePatientResponse, error)
}

func (m *mockPatientUsecase) FindByNationalIDOrPassportID(id string) (dto.PatientSearchResponse, error) {
	return m.findFunc(id)
}
func (m *mockPatientUsecase) CreatePatient(req dto.CreatePatientRequest, staffID uint, hospitalID uint) (*dto.CreatePatientResponse, error) {
	return m.createFunc(req, staffID, hospitalID)
}

func setupRouter(handler *PatientHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/patient/search/:id", handler.FindByNationalIDOrPassportID)
	r.POST("/patient", func(c *gin.Context) {
		c.Set("staff_id", uint(1))
		c.Set("hospital_id", uint(1))
		c.Next()
	}, handler.CreatePatient)
	return r
}

func TestHandler_SearchPatient_Success(t *testing.T) {
	uc := &mockPatientUsecase{
		findFunc: func(id string) (dto.PatientSearchResponse, error) {
			return dto.PatientSearchResponse{
				FirstNameTH: "สมชาย",
				PatientHN:   "OPD000001",
				NationalID:  id,
				DateOfBirth: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
				Gender:      "Male",
			}, nil
		},
	}
	handler := NewPatientHandler(uc)
	r := setupRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/patient/search/1234567890123", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.PatientSearchResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "สมชาย", resp.FirstNameTH)
	assert.Equal(t, "OPD000001", resp.PatientHN)
}

func TestHandler_SearchPatient_NotFound(t *testing.T) {
	uc := &mockPatientUsecase{
		findFunc: func(id string) (dto.PatientSearchResponse, error) {
			return dto.PatientSearchResponse{}, errors.New("record not found")
		},
	}
	handler := NewPatientHandler(uc)
	r := setupRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/patient/search/9999999999999", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestHandler_CreatePatient_Success(t *testing.T) {
	uc := &mockPatientUsecase{
		createFunc: func(req dto.CreatePatientRequest, staffID uint, hospitalID uint) (*dto.CreatePatientResponse, error) {
			return &dto.CreatePatientResponse{
				ID:          1,
				PatientHN:   "OPD000001",
				FirstNameTH: req.FirstNameTH,
				PatientType: "Outpatient",
			}, nil
		},
	}
	handler := NewPatientHandler(uc)
	r := setupRouter(handler)

	body, _ := json.Marshal(dto.CreatePatientRequest{
		FirstNameTH:   "สมชาย",
		LastNameTH:    "ใจดี",
		FirstNameEN:   "Somchai",
		LastNameEN:    "Jaidee",
		DateOfBirth:   "1990-01-15",
		NationalID:    "1234567890123",
		RefGenderID:   1,
		PatientTypeID: 2,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/patient", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp dto.CreatePatientResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "OPD000001", resp.PatientHN)
}

func TestHandler_CreatePatient_InvalidBody(t *testing.T) {
	uc := &mockPatientUsecase{}
	handler := NewPatientHandler(uc)
	r := setupRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/patient", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
