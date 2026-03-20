package route

import (
	"agnos-test/config"
	"agnos-test/middleware"
	authhandler "agnos-test/src/auth/auth_handler"
	authrepository "agnos-test/src/auth/auth_repository"
	authusecase "agnos-test/src/auth/auth_usecase"
	patienthandler "agnos-test/src/patient/patient_handler"
	patientrepository "agnos-test/src/patient/patient_repository"
	patientusecase "agnos-test/src/patient/patient_usecase"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(cfg *config.Config) *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "OK"})
	})

	hospitalRoutes := router.Group("/")
	hospitalRoutes.Use(middleware.HospitalMiddleware(cfg.DB))
	{
		// Auth
		authRepo := authrepository.NewAuthRepository(cfg.DB)
		authUsecase := authusecase.NewAuthUsecase(authRepo, cfg.JWTSecret)
		authHandler := authhandler.NewAuthHandler(authUsecase)

		hospitalRoutes.POST("/login", authHandler.Login)

		// Patient (protected)
		patientRepo := patientrepository.NewPatientRepository(cfg.DB)
		patientUsecase := patientusecase.NewPatientUsecase(patientRepo)
		patientHandler := patienthandler.NewPatientHandler(patientUsecase)

		authorized := hospitalRoutes.Group("/")
		authorized.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			authorized.GET("/patient/search/:id", patientHandler.FindByNationalIDOrPassportID)
			authorized.POST("/patient", patientHandler.CreatePatient)
		}
	}

	return router
}
