package usecase

import (
	"context"
	"dealls-dating/internal/entity"
	"dealls-dating/internal/model"
	"dealls-dating/internal/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SwipeUsecase struct {
	DB                    *gorm.DB
	Log                   *logrus.Logger
	matchRepository       *repository.MatchRepository
	swipeRepository       *repository.SwipeRepository
	userProfileRepository *repository.UserProfileRepository
	Validate              *validator.Validate
}

func NewSwipeUsecase(
	db *gorm.DB,
	log *logrus.Logger,
	matchRepository *repository.MatchRepository,
	swipeRepository *repository.SwipeRepository,
	validate *validator.Validate,
) *SwipeUsecase {
	return &SwipeUsecase{
		DB:              db,
		Log:             log,
		matchRepository: matchRepository,
		swipeRepository: swipeRepository,
		Validate:        validate,
	}
}

func (u *SwipeUsecase) Dislike(ctx context.Context, userId uint, swipeUserId uint) error {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	totalSwiped, err := u.swipeRepository.CountByUserIdAndDate(tx, userId, time.Now())
	if err != nil {
		return fiber.ErrInternalServerError
	} else if totalSwiped >= 10 {
		return fiber.ErrTooManyRequests
	}

	swiped, err := u.swipeRepository.CountByUserIdAndSwipeUserId(tx, userId, swipeUserId)
	if err != nil {
		return fiber.ErrInternalServerError
	} else if swiped > 0 {
		return fiber.NewError(fiber.StatusConflict, "user already swiped")
	}

	swipe := new(entity.Swipe)
	swipe.UserId = userId
	swipe.SwipeUserId = swipeUserId
	swipe.SwipeLike = false

	if err := u.swipeRepository.Create(tx, swipe); err != nil {
		u.Log.Warnf("Failed create swipe : %+v", err)
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed commit transaction : %+v", err)
		return fiber.ErrInternalServerError
	}

	return nil
}

func (u *SwipeUsecase) Find(ctx context.Context, userId uint) (*model.SwipeUserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	userProfile := new(entity.UserProfile)
	if err := u.userProfileRepository.FindByUserId(tx, userProfile, userId); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fiber.NewError(fiber.StatusPreconditionFailed, "user profile not completed")
		}
		u.Log.Warnf("Failed find user profile : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	swipeUserProfile := new(entity.UserProfile)
	if err := u.swipeRepository.FindAvailableSwipeUser(tx, swipeUserProfile, userId); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fiber.ErrNotFound
		}
		u.Log.Warnf("Failed find available swipe user : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return &model.SwipeUserResponse{
		UserId:     swipeUserProfile.UserId,
		Gender:     string(swipeUserProfile.Gender),
		Name:       swipeUserProfile.Name,
		PictureURL: swipeUserProfile.PictureURL,
	}, nil
}

func (u *SwipeUsecase) Like(ctx context.Context, userId uint, swipeUserId uint) (*model.SwipeLikeResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	totalSwiped, err := u.swipeRepository.CountByUserIdAndDate(tx, userId, time.Now())
	if err != nil {
		return nil, fiber.ErrInternalServerError
	} else if totalSwiped >= 10 {
		return nil, fiber.ErrTooManyRequests
	}

	swiped, err := u.swipeRepository.CountByUserIdAndSwipeUserId(tx, userId, swipeUserId)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	} else if swiped > 0 {
		return nil, fiber.NewError(fiber.StatusConflict, "user already swiped")
	}

	swipe := new(entity.Swipe)
	swipe.UserId = userId
	swipe.SwipeUserId = swipeUserId
	swipe.SwipeLike = true

	if err := u.swipeRepository.Create(tx, swipe); err != nil {
		u.Log.Warnf("Failed create swipe : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	oppositeTotal, err := u.swipeRepository.CountByUserIdAndSwipeUserIdAndLike(tx, swipeUserId, userId)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	match := new(entity.Match)
	if oppositeTotal > 0 {
		match.UserId1 = userId
		match.UserId2 = swipeUserId
		if err := u.matchRepository.Create(tx, match); err != nil {
			u.Log.Warnf("Failed create match : %+v", err)
			return nil, fiber.ErrInternalServerError
		}
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return &model.SwipeLikeResponse{
		Match: match.ID > 0,
	}, nil
}
