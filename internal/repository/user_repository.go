package repository

import (
	"dealls-dating/internal/entity"

	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[entity.User]
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (u *UserRepository) CountByEmail(db *gorm.DB, email string) (total int64, err error) {
	err = db.Model(entity.User{}).Where("email = ?", email).Count(&total).Error
	return total, err
}

func (u *UserRepository) FindByEmail(db *gorm.DB, user *entity.User, email string) error {
	return db.Where("email = ?", email).Take(user).Error
}

func (u *UserRepository) FindByToken(db *gorm.DB, user *entity.User, token string) error {
	return db.Where("token = ?", token).Take(user).Error
}
