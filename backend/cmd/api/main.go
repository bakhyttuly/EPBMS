package main

import (
	"log"
	"os"

	"epbms/config"
	"epbms/internal/handler"
	appMiddleware "epbms/internal/middleware"
	"epbms/internal/repository"
	"epbms/internal/service"
	"epbms/pkg/logger"
	"epbms/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialise structured logger.
	appLog := logger.New()

	// Connect to the database and run migrations.
	db, err := config.NewDatabase()
	if err != nil {
		log.Fatalf("database initialisation failed: %v", err)
	}

	// --- Repository layer ---
	userRepo := repository.NewUserRepository(db)
	performerRepo := repository.NewPerformerRepository(db)
	bookingRepo := repository.NewBookingRepository(db)

	// --- Service layer ---
	authSvc := service.NewAuthService(userRepo, performerRepo, appLog)
	performerSvc := service.NewPerformerService(performerRepo, appLog)
	bookingSvc := service.NewBookingService(bookingRepo, performerRepo, appLog)

	// --- Handler layer ---
	authH := handler.NewAuthHandler(authSvc)
	performerH := handler.NewPerformerHandler(performerSvc)
	bookingH := handler.NewBookingHandler(bookingSvc)
	adminH := handler.NewAdminHandler(userRepo, bookingRepo, performerRepo)

	// --- HTTP server ---
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(appMiddleware.RequestLogger(appLog))
	r.Use(appMiddleware.RateLimiter(10, 30)) // 10 req/s, burst of 30

	routes.SetupRoutes(r, authH, performerH, bookingH, adminH)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	appLog.Info("EPBMS server starting", "port", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
