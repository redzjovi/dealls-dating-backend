package route

import (
	"dealls-dating/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                   *fiber.App
	AuthController        *http.AuthController
	AuthMiddleware        fiber.Handler
	SwipeController       *http.SwipeController
	UserPremiumController *http.UserPremiumController
	UserProfileController *http.UserProfileController
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.Post("/api/auth/login", c.AuthController.Login)
	c.App.Post("/api/auth/sign-up", c.AuthController.SignUp)
}

func (c *RouteConfig) SetupAuthRoute() {
	c.App.Use(c.AuthMiddleware)
	c.App.Delete("/api/auth/user", c.AuthController.Logout)
	c.App.Get("/api/auth/user/premium", c.UserPremiumController.List)
	c.App.Post("/api/auth/user/premium", c.UserPremiumController.Trial)
	c.App.Get("/api/auth/user/profile", c.UserProfileController.Find)
	c.App.Put("/api/auth/user/profile", c.UserProfileController.Update)
	c.App.Get("/api/auth/user/swipe", c.SwipeController.Find)
	c.App.Post("/api/auth/user/swipe/:userId", c.SwipeController.Like)
	c.App.Delete("/api/auth/user/swipe/:userId", c.SwipeController.Dislike)
}
