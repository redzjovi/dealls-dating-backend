package repository

import (
	"dealls-dating/internal/entity"
	"time"

	"gorm.io/gorm"
)

type UserPremiumRepository struct {
	Repository[entity.UserPremium]
}

func NewUserPremiumRepository() *UserPremiumRepository {
	return &UserPremiumRepository{}
}

func (r *UserPremiumRepository) ListByUserIdAndNow(db *gorm.DB, userId uint) (res []entity.UserPremium, err error) {
	err = db.
		Where("user_id = ?", userId).
		Where("start_at <= ?", time.Now()).
		Where("end_at >= ?", time.Now()).
		Find(&res).
		Error
	return res, err
}

func (r *UserPremiumRepository) CountByUserIdAndNow(db *gorm.DB, userId uint) (total int64, err error) {
	err = db.
		Model(entity.UserPremium{}).
		Where("user_id = ?", userId).
		Where("start_at <= ?", time.Now()).
		Where("end_at >= ?", time.Now()).
		Count(&total).
		Error
	return total, err
}
