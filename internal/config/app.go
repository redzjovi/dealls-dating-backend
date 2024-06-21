package config

import (
	"dealls-dating/internal/delivery/http"
	"dealls-dating/internal/delivery/http/middleware"
	"dealls-dating/internal/delivery/http/route"
	"dealls-dating/internal/repository"
	"dealls-dating/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	App      *fiber.App
	DB       *gorm.DB
	Log      *logrus.Logger
	Validate *validator.Validate
}

func Bootstrap(config *BootstrapConfig) {
	// setup repositories
	userRepository := repository.NewUserRepository()

	// setup use cases
	authUsecase := usecase.NewAuthUsecase(config.DB, config.Log, userRepository, config.Validate)

	// setup controller
	authController := http.NewAuthController(authUsecase, config.Log)

	// setup middleware
	authMiddleware := middleware.NewAuth(authUsecase)

	routeConfig := route.RouteConfig{
		App:            config.App,
		AuthController: authController,
		AuthMiddleware: authMiddleware,
	}
	routeConfig.Setup()
}
