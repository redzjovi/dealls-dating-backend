package http

import (
	"dealls-dating/internal/delivery/http/middleware"
	"dealls-dating/internal/model"
	"dealls-dating/internal/usecase"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserPremiumController struct {
	Log                *logrus.Logger
	UserPremiumUsecase *usecase.UserPremiumUsecase
}

func NewUserPremiumController(
	log *logrus.Logger,
	userUsecase *usecase.UserPremiumUsecase,
) *UserPremiumController {
	return &UserPremiumController{
		Log:                log,
		UserPremiumUsecase: userUsecase,
	}
}

func (c *UserPremiumController) List(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	response, err := c.UserPremiumUsecase.ListByUserId(ctx.UserContext(), auth.ID)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[[]model.UserPremiumResponse]{Data: response})
}

func (c *UserPremiumController) Trial(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	if err := c.UserPremiumUsecase.Trial(ctx.UserContext(), auth.ID); err != nil {
		return err
	}

	return ctx.SendStatus(http.StatusNoContent)
}
