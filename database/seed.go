package database

import (
	"agnos-test/model"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	seedRefGender(db)
	seedRefPatientType(db)
	seedHospitals(db)
	seedStaff(db)
	log.Println("database seed completed")
}

func seedRefGender(db *gorm.DB) {
	genders := []model.RefGender{
		{ID: 1, Code: "M", Name: "Male"},
		{ID: 2, Code: "F", Name: "Female"},
		{ID: 3, Code: "O", Name: "Other"},
	}
	for _, g := range genders {
		db.FirstOrCreate(&g, model.RefGender{ID: g.ID})
	}
}

func seedRefPatientType(db *gorm.DB) {
	types := []model.RefPatientType{
		{ID: 1, Code: "IPD", Name: "Inpatient"},
		{ID: 2, Code: "OPD", Name: "Outpatient"},
		{ID: 3, Code: "EMS", Name: "Emergency"},
	}
	for _, t := range types {
		db.FirstOrCreate(&t, model.RefPatientType{ID: t.ID})
	}
}

func seedHospitals(db *gorm.DB) {
	hospitals := []model.Hospital{
		{ID: 1, Name: "โรงพยาบาล A", Slug: "hospital-a"},
		{ID: 2, Name: "โรงพยาบาล B", Slug: "hospital-b"},
	}
	for _, h := range hospitals {
		db.FirstOrCreate(&h, model.Hospital{ID: h.ID})
	}
}

func seedStaff(db *gorm.DB) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("failed to hash password for seed: %v", err)
		return
	}

	staffs := []model.Staff{
		{ID: 1, Username: "admin", Password: string(hashedPassword), HospitalID: 1},
		{ID: 2, Username: "admin", Password: string(hashedPassword), HospitalID: 2},
	}
	for _, s := range staffs {
		db.FirstOrCreate(&s, model.Staff{ID: s.ID})
	}
}
