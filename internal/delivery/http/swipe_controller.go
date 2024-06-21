package http

import (
	"dealls-dating/internal/delivery/http/middleware"
	"dealls-dating/internal/model"
	"dealls-dating/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type SwipeController struct {
	Log          *logrus.Logger
	SwipeUsecase *usecase.SwipeUsecase
}

func NewSwipeController(
	log *logrus.Logger,
	swipeUsecase *usecase.SwipeUsecase,
) *SwipeController {
	return &SwipeController{
		Log:          log,
		SwipeUsecase: swipeUsecase,
	}
}

func (c *SwipeController) Dislike(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	userIdUint64, err := strconv.ParseUint(ctx.Params("userId"), 10, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}

	swipeUserId := uint(userIdUint64)

	if err := c.SwipeUsecase.Dislike(ctx.UserContext(), auth.ID, swipeUserId); err != nil {
		return err
	}

	return ctx.SendStatus(http.StatusNoContent)
}

func (c *SwipeController) Find(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	response, err := c.SwipeUsecase.Find(ctx.UserContext(), auth.ID)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[*model.SwipeUserResponse]{Data: response})
}

func (c *SwipeController) Like(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	userIdUint64, err := strconv.ParseUint(ctx.Params("userId"), 10, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}

	swipeUserId := uint(userIdUint64)

	response, err := c.SwipeUsecase.Like(ctx.UserContext(), auth.ID, swipeUserId)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[*model.SwipeLikeResponse]{Data: response})
}
