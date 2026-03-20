package patientusecase

import (
	"agnos-test/dto"
	"agnos-test/model"
	patientrepository "agnos-test/src/patient/patient_repository"
	"errors"
	"time"
)

type PatientUsecase interface {
	FindByNationalIDOrPassportID(id string) (dto.PatientSearchResponse, error)
	CreatePatient(req dto.CreatePatientRequest, staffID uint, hospitalID uint) (*dto.CreatePatientResponse, error)
}

type patientUsecase struct {
	patientRepository patientrepository.PatientRepository
}

func NewPatientUsecase(patientRepository patientrepository.PatientRepository) PatientUsecase {
	return &patientUsecase{patientRepository: patientRepository}
}

func (u *patientUsecase) FindByNationalIDOrPassportID(id string) (dto.PatientSearchResponse, error) {
	patient, err := u.patientRepository.FindByNationalIDOrPassportID(id)
	if err != nil {
		return dto.PatientSearchResponse{}, err
	}
	return dto.PatientSearchResponse{
		FirstNameTH:  patient.FirstNameTH,
		MiddleNameTH: patient.MiddleNameTH,
		LastNameTH:   patient.LastNameTH,
		FirstNameEN:  patient.FirstNameEN,
		MiddleNameEN: patient.MiddleNameEN,
		LastNameEN:   patient.LastNameEN,
		DateOfBirth:  patient.DateOfBirth,
		PatientHN:    patient.PatientHN,
		NationalID:   patient.NationalID,
		PassportID:   patient.PassportID,
		PhoneNumber:  patient.PhoneNumber,
		Email:        patient.Email,
		Gender:       patient.Gender.Name,
	}, nil
}

func (u *patientUsecase) CreatePatient(req dto.CreatePatientRequest, staffID uint, hospitalID uint) (*dto.CreatePatientResponse, error) {
	if req.NationalID == "" && req.PassportID == "" {
		return nil, errors.New("national_id or passport_id is required")
	}

	patientType, err := u.patientRepository.GetPatientType(req.PatientTypeID)
	if err != nil {
		return nil, errors.New("invalid patient_type_id")
	}

	hospital, err := u.patientRepository.GetHospital(hospitalID)
	if err != nil {
		return nil, errors.New("invalid hospital_id")
	}

	gender, err := u.patientRepository.GetGender(req.RefGenderID)
	if err != nil {
		return nil, errors.New("invalid ref_gender_id")
	}

	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return nil, errors.New("invalid date_of_birth format, expected YYYY-MM-DD")
	}

	patient := &model.Patient{
		FirstNameTH:   req.FirstNameTH,
		MiddleNameTH:  req.MiddleNameTH,
		LastNameTH:    req.LastNameTH,
		FirstNameEN:   req.FirstNameEN,
		MiddleNameEN:  req.MiddleNameEN,
		LastNameEN:    req.LastNameEN,
		DateOfBirth:   dob,
		NationalID:    req.NationalID,
		PassportID:    req.PassportID,
		PhoneNumber:   req.PhoneNumber,
		Email:         req.Email,
		RefGenderID:   req.RefGenderID,
		PatientTypeID: req.PatientTypeID,
		HospitalID:    hospitalID,
		CreatedBy:     staffID,
		UpdatedBy:     staffID,
	}

	if err := u.patientRepository.CreatePatient(patient, patientType.Code); err != nil {
		return nil, err
	}

	return &dto.CreatePatientResponse{
		ID:           patient.ID,
		PatientHN:    patient.PatientHN,
		FirstNameTH:  patient.FirstNameTH,
		MiddleNameTH: patient.MiddleNameTH,
		LastNameTH:   patient.LastNameTH,
		FirstNameEN:  patient.FirstNameEN,
		MiddleNameEN: patient.MiddleNameEN,
		LastNameEN:   patient.LastNameEN,
		DateOfBirth:  patient.DateOfBirth.Format("2006-01-02"),
		NationalID:   patient.NationalID,
		PassportID:   patient.PassportID,
		PhoneNumber:  patient.PhoneNumber,
		Email:        patient.Email,
		Gender:       gender.Name,
		PatientType:  patientType.Name,
		HospitalName: hospital.Name,
	}, nil
}
