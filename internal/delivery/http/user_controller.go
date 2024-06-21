package http

import (
	"dealls-dating/internal/delivery/http/middleware"
	"dealls-dating/internal/model"
	"dealls-dating/internal/usecase"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserProfileController struct {
	Log                *logrus.Logger
	UserProfileUsecase *usecase.UserProfileUsecase
}

func NewUserProfileController(
	log *logrus.Logger,
	userUsecase *usecase.UserProfileUsecase,
) *UserProfileController {
	return &UserProfileController{
		Log:                log,
		UserProfileUsecase: userUsecase,
	}
}

func (c *UserProfileController) Find(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	response, err := c.UserProfileUsecase.Find(ctx.UserContext(), auth.ID)
	if err != nil {
		c.Log.Warnf("Failed to find user profile : %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserProfileResponse]{Data: response})
}

func (c *UserProfileController) Update(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := new(model.UpdateUserProfileRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	if err := c.UserProfileUsecase.Update(ctx.UserContext(), auth.ID, request); err != nil {
		c.Log.Warnf("Failed to update user profile : %+v", err)
		return err
	}

	return ctx.SendStatus(http.StatusNoContent)
}
