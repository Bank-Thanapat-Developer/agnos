package database

import (
	"agnos-test/model"
	"log"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.RefGender{},
		&model.RefPatientType{},
		&model.Hospital{},
		&model.Staff{},
		&model.Patient{},
		&model.PatientCodeSequences{},
	)
	if err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	log.Println("database migration completed")
}
