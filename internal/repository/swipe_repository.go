package repository

import (
	"dealls-dating/internal/entity"
	"time"

	"gorm.io/gorm"
)

type SwipeRepository struct {
	Repository[entity.Swipe]
}

func NewSwipeRepository() *SwipeRepository {
	return &SwipeRepository{}
}

func (r *SwipeRepository) CountByUserIdAndDate(db *gorm.DB, userId uint, date time.Time) (total int64, err error) {
	err = db.
		Model(entity.Swipe{}).
		Where("user_id = ?", userId).
		Where("DATE(created_at) = ?", date.Format("2006-01-02")).
		Count(&total).
		Error
	return total, err
}

func (r *SwipeRepository) CountByUserIdAndSwipeUserId(db *gorm.DB, userId uint, swipeUserId uint) (total int64, err error) {
	err = db.
		Model(entity.Swipe{}).
		Where("user_id = ?", userId).
		Where("swipe_user_id = ?", swipeUserId).
		Count(&total).
		Error
	return total, err
}

func (r *SwipeRepository) CountByUserIdAndSwipeUserIdAndLike(db *gorm.DB, userId uint, swipeUserId uint) (total int64, err error) {
	err = db.
		Model(entity.Swipe{}).
		Where("user_id = ?", userId).
		Where("swipe_user_id = ?", swipeUserId).
		Where("swipe_like = ?", true).
		Count(&total).
		Error
	return total, err
}

func (r *SwipeRepository) FindAvailableSwipeUser(db *gorm.DB, userProfile *entity.UserProfile, userId uint) (err error) {
	return db.
		Table("user_profiles").
		Where("user_id NOT IN (?)", db.Select("swipe_user_id").Table("swipes").Where("user_id = ?", userId)).
		Where("user_id <> ?", userId).
		Where("gender <> (?)", db.Select("gender").Table("user_profiles").Where("user_id = ?", userId)).
		Order("RANDOM()").
		Limit(1).
		Take(userProfile).Error
}
