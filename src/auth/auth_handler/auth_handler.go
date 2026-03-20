package authhandler

import (
	"agnos-test/dto"
	authusecase "agnos-test/src/auth/auth_usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase authusecase.AuthUsecase
}

func NewAuthHandler(authUsecase authusecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalID, _ := c.Get("hospital_id")

	resp, err := h.authUsecase.Login(req, hospitalID.(uint))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
