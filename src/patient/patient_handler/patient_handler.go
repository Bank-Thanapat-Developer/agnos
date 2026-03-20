package patienthandler

import (
	"agnos-test/dto"
	patientusecase "agnos-test/src/patient/patient_usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PatientHandler struct {
	patientUsecase patientusecase.PatientUsecase
}

func NewPatientHandler(patientUsecase patientusecase.PatientUsecase) *PatientHandler {
	return &PatientHandler{patientUsecase: patientUsecase}
}

func (h *PatientHandler) FindByNationalIDOrPassportID(c *gin.Context) {
	id := c.Param("id")
	patient, err := h.patientUsecase.FindByNationalIDOrPassportID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, patient)
}

func (h *PatientHandler) CreatePatient(c *gin.Context) {
	var req dto.CreatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	staffID, _ := c.Get("staff_id")
	hospitalID, _ := c.Get("hospital_id")

	resp, err := h.patientUsecase.CreatePatient(req, staffID.(uint), hospitalID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}
