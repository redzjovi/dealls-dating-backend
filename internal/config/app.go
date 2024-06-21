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
	matchRepository := repository.NewMatchRepository()
	swipeRepository := repository.NewSwipeRepository()
	userPremiumRepository := repository.NewUserPremiumRepository()
	userProfileRepository := repository.NewUserProfileRepository()
	userRepository := repository.NewUserRepository()

	// setup use cases
	authUsecase := usecase.NewAuthUsecase(config.DB, config.Log, userRepository, config.Validate)
	swipeUsecase := usecase.NewSwipeUsecase(config.DB, config.Log, matchRepository, swipeRepository, userPremiumRepository, userProfileRepository, config.Validate)
	userPremiumUsecase := usecase.NewUserPremiumUsecase(config.DB, config.Log, userPremiumRepository, config.Validate)
	userProfileUsecase := usecase.NewUserProfileUsecase(config.DB, config.Log, userProfileRepository, config.Validate)

	// setup controller
	authController := http.NewAuthController(authUsecase, config.Log)
	swipeController := http.NewSwipeController(config.Log, swipeUsecase)
	userPremiumController := http.NewUserPremiumController(config.Log, userPremiumUsecase)
	userProfileController := http.NewUserProfileController(config.Log, userProfileUsecase)

	// setup middleware
	authMiddleware := middleware.NewAuth(authUsecase)

	routeConfig := route.RouteConfig{
		App:                   config.App,
		AuthController:        authController,
		AuthMiddleware:        authMiddleware,
		SwipeController:       swipeController,
		UserPremiumController: userPremiumController,
		UserProfileController: userProfileController,
	}
	routeConfig.Setup()
}
