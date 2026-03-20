package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dsn string) *gorm.DB {

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// if err := db.AutoMigrate(&model.Patient{}, &model.Staff{}); err != nil {
	// 	log.Fatalf("failed to run migrations: %v", err)
	// }

	return db
}
