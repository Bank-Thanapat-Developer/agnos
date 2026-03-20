package patientusecase

import (
	"agnos-test/dto"
	"agnos-test/model"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockPatientRepository struct {
	findByIDFunc      func(id string) (*model.PatientHistory, error)
	createPatientFunc func(patient *model.Patient, typeCode string) error
	getPatientType    func(id uint) (*model.RefPatientType, error)
	getHospital       func(id uint) (*model.Hospital, error)
	getGender         func(id uint) (*model.RefGender, error)
}

func (m *mockPatientRepository) FindByNationalIDOrPassportID(id string) (*model.PatientHistory, error) {
	return m.findByIDFunc(id)
}
func (m *mockPatientRepository) CreatePatient(patient *model.Patient, typeCode string) error {
	return m.createPatientFunc(patient, typeCode)
}
func (m *mockPatientRepository) GetPatientType(id uint) (*model.RefPatientType, error) {
	return m.getPatientType(id)
}
func (m *mockPatientRepository) GetHospital(id uint) (*model.Hospital, error) {
	return m.getHospital(id)
}
func (m *mockPatientRepository) GetGender(id uint) (*model.RefGender, error) {
	return m.getGender(id)
}

func defaultMockRepo() *mockPatientRepository {
	return &mockPatientRepository{
		getPatientType: func(id uint) (*model.RefPatientType, error) {
			types := map[uint]*model.RefPatientType{
				1: {ID: 1, Code: "IPD", Name: "Inpatient"},
				2: {ID: 2, Code: "OPD", Name: "Outpatient"},
				3: {ID: 3, Code: "EMS", Name: "Emergency"},
			}
			if t, ok := types[id]; ok {
				return t, nil
			}
			return nil, errors.New("not found")
		},
		getHospital: func(id uint) (*model.Hospital, error) {
			if id == 1 {
				return &model.Hospital{ID: 1, Name: "Hospital A", Slug: "hospital-a"}, nil
			}
			return nil, errors.New("not found")
		},
		getGender: func(id uint) (*model.RefGender, error) {
			genders := map[uint]*model.RefGender{
				1: {ID: 1, Code: "M", Name: "Male"},
				2: {ID: 2, Code: "F", Name: "Female"},
			}
			if g, ok := genders[id]; ok {
				return g, nil
			}
			return nil, errors.New("not found")
		},
		createPatientFunc: func(patient *model.Patient, typeCode string) error {
			patient.ID = 1
			patient.PatientHN = typeCode + "000001"
			return nil
		},
	}
}

func validCreateRequest() dto.CreatePatientRequest {
	return dto.CreatePatientRequest{
		FirstNameTH:   "สมชาย",
		LastNameTH:    "ใจดี",
		FirstNameEN:   "Somchai",
		LastNameEN:    "Jaidee",
		DateOfBirth:   "1990-01-15",
		NationalID:    "1234567890123",
		PhoneNumber:   "0812345678",
		RefGenderID:   1,
		PatientTypeID: 2,
	}
}

func TestCreatePatient_OPD_Success(t *testing.T) {
	repo := defaultMockRepo()
	uc := NewPatientUsecase(repo)

	resp, err := uc.CreatePatient(validCreateRequest(), 1, 1)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "OPD000001", resp.PatientHN)
	assert.Equal(t, "สมชาย", resp.FirstNameTH)
	assert.Equal(t, "Outpatient", resp.PatientType)
	assert.Equal(t, "Male", resp.Gender)
	assert.Equal(t, "Hospital A", resp.HospitalName)
}

func TestCreatePatient_IPD_Success(t *testing.T) {
	repo := defaultMockRepo()
	repo.createPatientFunc = func(patient *model.Patient, typeCode string) error {
		patient.ID = 2
		patient.PatientHN = typeCode + "000001"
		return nil
	}
	uc := NewPatientUsecase(repo)

	req := validCreateRequest()
	req.PatientTypeID = 1
	resp, err := uc.CreatePatient(req, 1, 1)

	assert.NoError(t, err)
	assert.Equal(t, "IPD000001", resp.PatientHN)
	assert.Equal(t, "Inpatient", resp.PatientType)
}

func TestCreatePatient_EMS_Success(t *testing.T) {
	repo := defaultMockRepo()
	repo.createPatientFunc = func(patient *model.Patient, typeCode string) error {
		patient.ID = 3
		patient.PatientHN = typeCode + "000001"
		return nil
	}
	uc := NewPatientUsecase(repo)

	req := validCreateRequest()
	req.PatientTypeID = 3
	req.NationalID = ""
	req.PassportID = "AB1234567"
	resp, err := uc.CreatePatient(req, 1, 1)

	assert.NoError(t, err)
	assert.Equal(t, "EMS000001", resp.PatientHN)
	assert.Equal(t, "Emergency", resp.PatientType)
}

func TestCreatePatient_MissingNationalIDAndPassport(t *testing.T) {
	repo := defaultMockRepo()
	uc := NewPatientUsecase(repo)

	req := validCreateRequest()
	req.NationalID = ""
	req.PassportID = ""
	resp, err := uc.CreatePatient(req, 1, 1)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "national_id or passport_id is required", err.Error())
}

func TestCreatePatient_InvalidPatientType(t *testing.T) {
	repo := defaultMockRepo()
	uc := NewPatientUsecase(repo)

	req := validCreateRequest()
	req.PatientTypeID = 99
	resp, err := uc.CreatePatient(req, 1, 1)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "invalid patient_type_id", err.Error())
}

func TestCreatePatient_InvalidHospital(t *testing.T) {
	repo := defaultMockRepo()
	uc := NewPatientUsecase(repo)

	req := validCreateRequest()
	resp, err := uc.CreatePatient(req, 1, 99)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "invalid hospital_id", err.Error())
}

func TestCreatePatient_InvalidGender(t *testing.T) {
	repo := defaultMockRepo()
	uc := NewPatientUsecase(repo)

	req := validCreateRequest()
	req.RefGenderID = 99
	resp, err := uc.CreatePatient(req, 1, 1)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "invalid ref_gender_id", err.Error())
}

func TestCreatePatient_InvalidDateFormat(t *testing.T) {
	repo := defaultMockRepo()
	uc := NewPatientUsecase(repo)

	req := validCreateRequest()
	req.DateOfBirth = "15-01-1990"
	resp, err := uc.CreatePatient(req, 1, 1)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "invalid date_of_birth format, expected YYYY-MM-DD", err.Error())
}

func TestCreatePatient_RepositoryError(t *testing.T) {
	repo := defaultMockRepo()
	repo.createPatientFunc = func(patient *model.Patient, typeCode string) error {
		return errors.New("db error")
	}
	uc := NewPatientUsecase(repo)

	resp, err := uc.CreatePatient(validCreateRequest(), 1, 1)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestFindByNationalIDOrPassportID_Success(t *testing.T) {
	repo := defaultMockRepo()
	repo.findByIDFunc = func(id string) (*model.PatientHistory, error) {
		return &model.PatientHistory{
			ID:          1,
			FirstNameTH: "สมชาย",
			LastNameTH:  "ใจดี",
			FirstNameEN: "Somchai",
			LastNameEN:  "Jaidee",
			DateOfBirth: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
			PatientHN:   "OPD000001",
			NationalID:  "1234567890123",
			Gender:      model.RefGender{Name: "Male"},
		}, nil
	}
	uc := NewPatientUsecase(repo)

	resp, err := uc.FindByNationalIDOrPassportID("1234567890123")

	assert.NoError(t, err)
	assert.Equal(t, "สมชาย", resp.FirstNameTH)
	assert.Equal(t, "OPD000001", resp.PatientHN)
	assert.Equal(t, "Male", resp.Gender)
}

func TestFindByNationalIDOrPassportID_NotFound(t *testing.T) {
	repo := defaultMockRepo()
	repo.findByIDFunc = func(id string) (*model.PatientHistory, error) {
		return nil, errors.New("record not found")
	}
	uc := NewPatientUsecase(repo)

	_, err := uc.FindByNationalIDOrPassportID("9999999999999")

	assert.Error(t, err)
}
