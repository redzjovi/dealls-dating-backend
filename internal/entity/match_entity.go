package entity

import "time"

type Match struct {
	ID        uint      `gorm:"column:id;primaryKey"`
	UserId1   uint      `gorm:"column:user_id_1"`
	UserId2   uint      `gorm:"column:user_id_2"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}
