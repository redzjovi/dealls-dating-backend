package entity

import "time"

type UserPremium struct {
	ID        uint      `gorm:"column:id;primaryKey"`
	UserId    uint      `gorm:"column:user_id"`
	StartAt   time.Time `gorm:"column:start_at"`
	EndAt     time.Time `gorm:"column:end_at"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (u *UserPremium) TableName() string {
	return "user_premiums"
}
