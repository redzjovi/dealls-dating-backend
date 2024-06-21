package usecase

import (
	"context"
	"dealls-dating/internal/entity"
	"dealls-dating/internal/model"
	"dealls-dating/internal/model/converter"
	"dealls-dating/internal/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserPremiumUsecase struct {
	DB                    *gorm.DB
	Log                   *logrus.Logger
	UserPremiumRepository *repository.UserPremiumRepository
	Validate              *validator.Validate
}

func NewUserPremiumUsecase(
	db *gorm.DB,
	log *logrus.Logger,
	UserPremiumRepository *repository.UserPremiumRepository,
	validate *validator.Validate,
) *UserPremiumUsecase {
	return &UserPremiumUsecase{
		DB:                    db,
		Log:                   log,
		UserPremiumRepository: UserPremiumRepository,
		Validate:              validate,
	}
}

func (u *UserPremiumUsecase) ListByUserId(ctx context.Context, userId uint) ([]model.UserPremiumResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	userPremiums, err := u.UserPremiumRepository.ListByUserIdAndNow(tx, userId)
	if err != nil {
		u.Log.Warnf("Failed list user premium by user id : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	res := make([]model.UserPremiumResponse, 0)
	for _, userPremium := range userPremiums {
		res = append(res, *converter.UserPremiumToResponse(&userPremium))
	}

	return res, nil
}

func (u *UserPremiumUsecase) Trial(ctx context.Context, userId uint) error {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	total, err := u.UserPremiumRepository.CountByUserIdAndNow(tx, userId)
	if err != nil {
		u.Log.Warnf("Failed count user premium by user id : %+v", err)
		return fiber.ErrInternalServerError
	} else if total > 0 {
		return fiber.ErrConflict
	}

	userPremium := new(entity.UserPremium)
	userPremium.UserId = userId
	userPremium.StartAt = time.Now()
	userPremium.EndAt = time.Now().Add(7 * 24 * time.Hour)

	if err := u.UserPremiumRepository.Create(tx, userPremium); err != nil {
		u.Log.Warnf("Failed create user premium : %+v", err)
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed commit transaction : %+v", err)
		return fiber.ErrInternalServerError
	}

	return nil
}
