package dto

import "time"

type PatientSearchResponse struct {
	FirstNameTH  string    `json:"first_name_th"`
	MiddleNameTH string    `json:"middle_name_th"`
	LastNameTH   string    `json:"last_name_th"`
	FirstNameEN  string    `json:"first_name_en"`
	MiddleNameEN string    `json:"middle_name_en"`
	LastNameEN   string    `json:"last_name_en"`
	DateOfBirth  time.Time `json:"date_of_birth"`
	PatientHN    string    `json:"patient_hn"`
	NationalID   string    `json:"national_id"`
	PassportID   string    `json:"passport_id"`
	PhoneNumber  string    `json:"phone_number"`
	Email        string    `json:"email"`
	Gender       string    `json:"gender"`
}

type CreatePatientRequest struct {
	FirstNameTH   string `json:"first_name_th" binding:"required"`
	MiddleNameTH  string `json:"middle_name_th"`
	LastNameTH    string `json:"last_name_th" binding:"required"`
	FirstNameEN   string `json:"first_name_en" binding:"required"`
	MiddleNameEN  string `json:"middle_name_en"`
	LastNameEN    string `json:"last_name_en" binding:"required"`
	DateOfBirth   string `json:"date_of_birth" binding:"required"`
	NationalID    string `json:"national_id"`
	PassportID    string `json:"passport_id"`
	PhoneNumber   string `json:"phone_number"`
	Email         string `json:"email"`
	RefGenderID   uint   `json:"ref_gender_id" binding:"required"`
	PatientTypeID uint   `json:"patient_type_id" binding:"required"`
}

type CreatePatientResponse struct {
	ID           uint   `json:"id"`
	PatientHN    string `json:"patient_hn"`
	FirstNameTH  string `json:"first_name_th"`
	MiddleNameTH string `json:"middle_name_th"`
	LastNameTH   string `json:"last_name_th"`
	FirstNameEN  string `json:"first_name_en"`
	MiddleNameEN string `json:"middle_name_en"`
	LastNameEN   string `json:"last_name_en"`
	DateOfBirth  string `json:"date_of_birth"`
	NationalID   string `json:"national_id"`
	PassportID   string `json:"passport_id"`
	PhoneNumber  string `json:"phone_number"`
	Email        string `json:"email"`
	Gender       string `json:"gender"`
	PatientType  string `json:"patient_type"`
	HospitalName string `json:"hospital_name"`
}
