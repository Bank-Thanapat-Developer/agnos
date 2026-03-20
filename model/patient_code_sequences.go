package model

import "time"

type PatientCodeSequences struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Code       string    `gorm:"type:varchar(5);not null;uniqueIndex:idx_hospital_code_seq" json:"code"`
	Sequence   int       `gorm:"not null;uniqueIndex:idx_hospital_code_seq" json:"sequence"`
	PatientID  uint      `gorm:"not null;index" json:"patient_id"`
	HospitalID uint      `gorm:"not null;uniqueIndex:idx_hospital_code_seq" json:"hospital_id"`
	CreatedAt  time.Time `json:"created_at"`
}

func (PatientCodeSequences) TableName() string {
	return "patient_code_sequences"
}
