package patientrepository

import (
	"agnos-test/model"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PatientRepository interface {
	FindByNationalIDOrPassportID(id string) (*model.PatientHistory, error)
	CreatePatient(patient *model.Patient, typeCode string) error
	GetPatientType(id uint) (*model.RefPatientType, error)
	GetHospital(id uint) (*model.Hospital, error)
	GetGender(id uint) (*model.RefGender, error)
}

type patientRepository struct {
	db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) PatientRepository {
	return &patientRepository{db: db}
}

func (r *patientRepository) FindByNationalIDOrPassportID(id string) (*model.PatientHistory, error) {
	var patient model.PatientHistory
	err := r.db.Model(&model.PatientHistory{}).
		Preload("Gender").
		Where("national_id = ? OR passport_id = ?", id, id).
		First(&patient).Error
	if err != nil {
		return nil, err
	}
	return &patient, nil
}

func (r *patientRepository) CreatePatient(patient *model.Patient, typeCode string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Lock the hospital row to serialize sequence generation per hospital
		var hospital model.Hospital
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&hospital, patient.HospitalID).Error; err != nil {
			return fmt.Errorf("hospital not found: %w", err)
		}

		var maxSeq int
		tx.Model(&model.PatientCodeSequences{}).
			Select("COALESCE(MAX(sequence), 0)").
			Where("hospital_id = ? AND code = ?", patient.HospitalID, typeCode).
			Scan(&maxSeq)

		nextSeq := maxSeq + 1
		patient.PatientHN = fmt.Sprintf("%s%06d", typeCode, nextSeq)

		if err := tx.Create(patient).Error; err != nil {
			return err
		}

		seq := model.PatientCodeSequences{
			Code:       typeCode,
			Sequence:   nextSeq,
			PatientID:  patient.ID,
			HospitalID: patient.HospitalID,
		}
		return tx.Create(&seq).Error
	})
}

func (r *patientRepository) GetPatientType(id uint) (*model.RefPatientType, error) {
	var pt model.RefPatientType
	if err := r.db.First(&pt, id).Error; err != nil {
		return nil, err
	}
	return &pt, nil
}

func (r *patientRepository) GetHospital(id uint) (*model.Hospital, error) {
	var h model.Hospital
	if err := r.db.First(&h, id).Error; err != nil {
		return nil, err
	}
	return &h, nil
}

func (r *patientRepository) GetGender(id uint) (*model.RefGender, error) {
	var g model.RefGender
	if err := r.db.First(&g, id).Error; err != nil {
		return nil, err
	}
	return &g, nil
}
