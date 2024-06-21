package http

import (
	"dealls-dating/internal/delivery/http/middleware"
	"dealls-dating/internal/model"
	"dealls-dating/internal/usecase"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AuthController struct {
	AuthUsecase *usecase.AuthUsecase
	Log         *logrus.Logger
}

func NewAuthController(
	authUsecase *usecase.AuthUsecase,
	log *logrus.Logger,
) *AuthController {
	return &AuthController{
		AuthUsecase: authUsecase,
		Log:         log,
	}
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	request := new(model.AuthLoginRequest)
	if err := ctx.BodyParser(request); err != nil {
		return fiber.ErrBadRequest
	}

	response, err := c.AuthUsecase.Login(ctx.UserContext(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AuthLoginResponse]{Data: response})
}

func (c *AuthController) Logout(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &model.AuthLogoutRequest{
		ID: auth.ID,
	}

	if err := c.AuthUsecase.Logout(ctx.UserContext(), request); err != nil {
		return err
	}

	return ctx.SendStatus(http.StatusNoContent)
}

func (c *AuthController) SignUp(ctx *fiber.Ctx) error {
	request := new(model.AuthSignUpRequest)
	if err := ctx.BodyParser(request); err != nil {
		return fiber.ErrBadRequest
	}

	if err := c.AuthUsecase.SignUp(ctx.UserContext(), request); err != nil {
		return err
	}

	return ctx.SendStatus(http.StatusNoContent)
}
