package usecase

import (
	"context"
	"dealls-dating/internal/entity"
	"dealls-dating/internal/model"
	"dealls-dating/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUsecase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	UserRepository *repository.UserRepository
	Validate       *validator.Validate
}

func NewAuthUsecase(
	db *gorm.DB,
	log *logrus.Logger,
	userRepository *repository.UserRepository,
	validate *validator.Validate,
) *AuthUsecase {
	return &AuthUsecase{
		DB:             db,
		Log:            log,
		UserRepository: userRepository,
		Validate:       validate,
	}
}

func (u *AuthUsecase) Login(ctx context.Context, request *model.AuthLoginRequest) (*model.AuthLoginResponse, error) {
	if err := u.Validate.Struct(request); err != nil {
		u.Log.Infof("Invalid request body  : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user := new(entity.User)
	if err := u.UserRepository.FindByEmail(tx, user, request.Email); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fiber.ErrUnauthorized
		}
		u.Log.Warnf("Failed find user by id : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return nil, fiber.ErrUnauthorized
	}

	user.Token = uuid.New().String()
	if err := u.UserRepository.Update(tx, user); err != nil {
		u.Log.Warnf("Failed save user : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return &model.AuthLoginResponse{
		Token: user.Token,
	}, nil
}

func (u *AuthUsecase) Logout(ctx context.Context, request *model.AuthLogoutRequest) error {
	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Invalid request body : %+v", err)
		return fiber.ErrBadRequest
	}

	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user := new(entity.User)
	if err := u.UserRepository.FindById(tx, user, request.ID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.ErrNotFound
		}
		u.Log.Warnf("Failed find user by id : %+v", err)
		return fiber.ErrInternalServerError
	}

	user.Token = ""

	if err := u.UserRepository.Update(tx, user); err != nil {
		u.Log.Warnf("Failed save user : %+v", err)
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed commit transaction : %+v", err)
		return fiber.ErrInternalServerError
	}

	return nil
}

func (u *AuthUsecase) SignUp(ctx context.Context, request *model.AuthSignUpRequest) error {
	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Invalid request body : %+v", err)
		return fiber.ErrBadRequest
	}

	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	total, err := u.UserRepository.CountByEmail(tx, request.Email)
	if err != nil {
		u.Log.Warnf("Failed count user by email : %+v", err)
		return fiber.ErrInternalServerError
	}

	if total > 0 {
		return fiber.ErrConflict
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		u.Log.Warnf("Failed to generate bcrype hash : %+v", err)
		return fiber.ErrInternalServerError
	}

	user := &entity.User{
		Email:    request.Email,
		Password: string(password),
	}

	if err := u.UserRepository.Create(tx, user); err != nil {
		u.Log.Warnf("Failed create user : %+v", err)
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed commit transaction : %+v", err)
		return fiber.ErrInternalServerError
	}

	return nil
}

func (u *AuthUsecase) Verify(ctx context.Context, request *model.AuthVerifyRequest) (*model.Auth, error) {
	err := u.Validate.Struct(request)
	if err != nil {
		u.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user := new(entity.User)
	if err := u.UserRepository.FindByToken(tx, user, request.Token); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fiber.ErrNotFound
		}
		u.Log.Warnf("Failed find user by token : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return &model.Auth{ID: user.ID}, nil
}
