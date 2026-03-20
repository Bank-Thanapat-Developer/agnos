package authrepository

import (
	"agnos-test/model"

	"gorm.io/gorm"
)

type AuthRepository interface {
	FindStaffByUsernameAndHospital(username string, hospitalID uint) (*model.Staff, error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) FindStaffByUsernameAndHospital(username string, hospitalID uint) (*model.Staff, error) {
	var staff model.Staff
	err := r.db.Preload("Hospital").
		Where("username = ? AND hospital_id = ?", username, hospitalID).
		First(&staff).Error
	if err != nil {
		return nil, err
	}
	return &staff, nil
}
