package config

import (
	"agnos-test/database"
	"fmt"
	"os"

	"gorm.io/gorm"
)

type Config struct {
	DB         *gorm.DB
	JWTSecret  string
	ServerPort string
}

type DatabaseConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func Load() *Config {
	dbConfig := DatabaseConfig{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "agnos-test"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Bangkok",
		dbConfig.DBHost, dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBName, dbConfig.DBPort, dbConfig.DBSSLMode)

	db := database.Connect(dsn)

	return &Config{
		DB:         db,
		JWTSecret:  getEnv("JWT_SECRET", "supersecretkey"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
