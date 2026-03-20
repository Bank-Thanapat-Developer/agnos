package model

import "time"

type Patient struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	FirstNameTH   string    `gorm:"type:varchar(255)" json:"first_name_th"`
	MiddleNameTH  string    `gorm:"type:varchar(255)" json:"middle_name_th"`
	LastNameTH    string    `gorm:"type:varchar(255)" json:"last_name_th"`
	FirstNameEN   string    `gorm:"type:varchar(255)" json:"first_name_en"`
	MiddleNameEN  string    `gorm:"type:varchar(255)" json:"middle_name_en"`
	LastNameEN    string    `gorm:"type:varchar(255)" json:"last_name_en"`
	DateOfBirth   time.Time `gorm:"type:date" json:"date_of_birth"`
	PatientHN     string    `gorm:"type:varchar(50)" json:"patient_hn"`
	NationalID    string    `gorm:"type:varchar(20);index" json:"national_id"`
	PassportID    string    `gorm:"type:varchar(50);index" json:"passport_id"`
	PhoneNumber   string    `gorm:"type:varchar(20)" json:"phone_number"`
	Email         string    `gorm:"type:varchar(255)" json:"email"`
	RefGenderID   uint      `gorm:"type:uint;index;not null" json:"ref_gender_id"`
	PatientTypeID uint      `gorm:"type:uint;index;not null" json:"patient_type_id"`
	HospitalID    uint      `gorm:"type:uint;index;not null" json:"hospital_id"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     uint      `gorm:"type:uint;index;not null" json:"created_by"`
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedBy     uint      `gorm:"type:uint;index;not null" json:"updated_by"`
}

func (b *Patient) TableName() string {
	return "patient"
}

type PatientHistory struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	FirstNameTH  string    `gorm:"type:varchar(255)" json:"first_name_th"`
	MiddleNameTH string    `gorm:"type:varchar(255)" json:"middle_name_th"`
	LastNameTH   string    `gorm:"type:varchar(255)" json:"last_name_th"`
	FirstNameEN  string    `gorm:"type:varchar(255)" json:"first_name_en"`
	MiddleNameEN string    `gorm:"type:varchar(255)" json:"middle_name_en"`
	LastNameEN   string    `gorm:"type:varchar(255)" json:"last_name_en"`
	DateOfBirth  time.Time `gorm:"type:date" json:"date_of_birth"`
	PatientHN    string    `gorm:"type:varchar(50)" json:"patient_hn"`
	NationalID   string    `gorm:"type:varchar(20);index" json:"national_id"`
	PassportID   string    `gorm:"type:varchar(50);index" json:"passport_id"`
	PhoneNumber  string    `gorm:"type:varchar(20)" json:"phone_number"`
	Email        string    `gorm:"type:varchar(255)" json:"email"`
	RefGenderID  uint      `gorm:"type:uint;index;not null" json:"ref_gender_id"`
	Gender       RefGender `gorm:"foreignKey:RefGenderID" json:"gender"`
}

func (b *PatientHistory) TableName() string {
	return "patient"
}
