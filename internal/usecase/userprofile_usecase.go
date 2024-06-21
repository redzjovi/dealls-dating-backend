package usecase

import (
	"context"
	"dealls-dating/internal/entity"
	"dealls-dating/internal/model"
	"dealls-dating/internal/model/converter"
	"dealls-dating/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserProfileUsecase struct {
	DB                    *gorm.DB
	Log                   *logrus.Logger
	UserProfileRepository *repository.UserProfileRepository
	Validate              *validator.Validate
}

func NewUserProfileUsecase(
	db *gorm.DB,
	log *logrus.Logger,
	UserProfileRepository *repository.UserProfileRepository,
	validate *validator.Validate,
) *UserProfileUsecase {
	return &UserProfileUsecase{
		DB:                    db,
		Log:                   log,
		UserProfileRepository: UserProfileRepository,
		Validate:              validate,
	}
}

func (c *UserProfileUsecase) Find(ctx context.Context, userId uint) (*model.UserProfileResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	userProfile := new(entity.UserProfile)
	if err := c.UserProfileRepository.FindByUserId(tx, userProfile, userId); err != nil {
		c.Log.Warnf("Failed find user by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	return converter.UserProfileToResponse(userProfile), nil
}

func (c *UserProfileUsecase) Update(ctx context.Context, userId uint, request *model.UpdateUserProfileRequest) error {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return fiber.ErrBadRequest
	}

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	userProfile := new(entity.UserProfile)
	if err := c.UserProfileRepository.FindByUserId(tx, userProfile, userId); err != nil && err != gorm.ErrRecordNotFound {
		c.Log.Warnf("Failed find user by id : %+v", err)
		return fiber.ErrNotFound
	}

	userProfile.UserId = userId
	userProfile.Gender = entity.UserProfileGender(request.Gender)
	userProfile.Name = request.Name
	userProfile.PictureURL = request.PictureURL

	if err := c.UserProfileRepository.Update(tx, userProfile); err != nil {
		c.Log.Warnf("Failed save user profile : %+v", err)
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return fiber.ErrInternalServerError
	}

	return nil
}
