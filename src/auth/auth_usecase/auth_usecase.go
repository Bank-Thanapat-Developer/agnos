package authusecase

import (
	"agnos-test/dto"
	authrepository "agnos-test/src/auth/auth_repository"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Login(req dto.LoginRequest, hospitalID uint) (*dto.LoginResponse, error)
}

type authUsecase struct {
	authRepository authrepository.AuthRepository
	jwtSecret      string
}

func NewAuthUsecase(authRepository authrepository.AuthRepository, jwtSecret string) AuthUsecase {
	return &authUsecase{
		authRepository: authRepository,
		jwtSecret:      jwtSecret,
	}
}

func (u *authUsecase) Login(req dto.LoginRequest, hospitalID uint) (*dto.LoginResponse, error) {
	staff, err := u.authRepository.FindStaffByUsernameAndHospital(req.Username, hospitalID)
	if err != nil {
		return nil, errors.New("invalid username or hospital")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(staff.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"staff_id":    staff.ID,
		"hospital_id": staff.HospitalID,
		"username":    staff.Username,
		"exp":         time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &dto.LoginResponse{
		Token:        tokenString,
		StaffID:      staff.ID,
		Username:     staff.Username,
		HospitalID:   staff.HospitalID,
		HospitalName: staff.Hospital.Name,
	}, nil
}
