package authusecase

import (
	"agnos-test/dto"
	"agnos-test/model"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type mockAuthRepository struct {
	findStaffFunc func(username string, hospitalID uint) (*model.Staff, error)
}

func (m *mockAuthRepository) FindStaffByUsernameAndHospital(username string, hospitalID uint) (*model.Staff, error) {
	return m.findStaffFunc(username, hospitalID)
}

func hashedPassword(plain string) string {
	h, _ := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(h)
}

func TestLogin_Success(t *testing.T) {
	repo := &mockAuthRepository{
		findStaffFunc: func(username string, hospitalID uint) (*model.Staff, error) {
			return &model.Staff{
				ID:         1,
				Username:   "admin",
				Password:   hashedPassword("password123"),
				HospitalID: 1,
				Hospital:   model.Hospital{ID: 1, Name: "Hospital A"},
			}, nil
		},
	}
	uc := NewAuthUsecase(repo, "test-secret")

	resp, err := uc.Login(dto.LoginRequest{Username: "admin", Password: "password123"}, 1)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, uint(1), resp.StaffID)
	assert.Equal(t, "admin", resp.Username)
	assert.Equal(t, uint(1), resp.HospitalID)
	assert.Equal(t, "Hospital A", resp.HospitalName)

	token, err := jwt.Parse(resp.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})
	assert.NoError(t, err)
	assert.True(t, token.Valid)

	claims := token.Claims.(jwt.MapClaims)
	assert.Equal(t, float64(1), claims["staff_id"])
	assert.Equal(t, float64(1), claims["hospital_id"])
	assert.Equal(t, "admin", claims["username"])
}

func TestLogin_InvalidUsername(t *testing.T) {
	repo := &mockAuthRepository{
		findStaffFunc: func(username string, hospitalID uint) (*model.Staff, error) {
			return nil, errors.New("not found")
		},
	}
	uc := NewAuthUsecase(repo, "test-secret")

	resp, err := uc.Login(dto.LoginRequest{Username: "wrong", Password: "password123"}, 1)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "invalid username or hospital", err.Error())
}

func TestLogin_InvalidPassword(t *testing.T) {
	repo := &mockAuthRepository{
		findStaffFunc: func(username string, hospitalID uint) (*model.Staff, error) {
			return &model.Staff{
				ID:       1,
				Username: "admin",
				Password: hashedPassword("password123"),
			}, nil
		},
	}
	uc := NewAuthUsecase(repo, "test-secret")

	resp, err := uc.Login(dto.LoginRequest{Username: "admin", Password: "wrongpass"}, 1)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "invalid password", err.Error())
}

func TestLogin_InvalidHospital(t *testing.T) {
	repo := &mockAuthRepository{
		findStaffFunc: func(username string, hospitalID uint) (*model.Staff, error) {
			return nil, errors.New("not found")
		},
	}
	uc := NewAuthUsecase(repo, "test-secret")

	resp, err := uc.Login(dto.LoginRequest{Username: "admin", Password: "password123"}, 999)

	assert.Error(t, err)
	assert.Nil(t, resp)
}
