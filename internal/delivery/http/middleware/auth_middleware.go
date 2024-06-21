package middleware

import (
	"dealls-dating/internal/model"
	"dealls-dating/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

func NewAuth(authUsecase *usecase.AuthUsecase) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		request := &model.AuthVerifyRequest{Token: ctx.Get("Authorization", "NOT_FOUND")}
		authUsecase.Log.Debugf("Authorization : %s", request.Token)

		auth, err := authUsecase.Verify(ctx.UserContext(), request)
		if err != nil {
			authUsecase.Log.Warnf("Failed find user by token : %+v", err)
			return fiber.ErrUnauthorized
		}

		authUsecase.Log.Debugf("User : %+v", auth.ID)
		ctx.Locals("auth", auth)
		return ctx.Next()
	}
}

func GetUser(ctx *fiber.Ctx) *model.Auth {
	return ctx.Locals("auth").(*model.Auth)
}
