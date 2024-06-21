package entity

import "time"

type Swipe struct {
	ID          uint      `gorm:"column:id;primaryKey"`
	UserId      uint      `gorm:"column:user_id"`
	SwipeUserId uint      `gorm:"column:swipe_user_id"`
	SwipeLike   bool      `gorm:"column:swipe_like"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}
