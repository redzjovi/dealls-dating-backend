package entity

import "time"

type User struct {
	ID        uint      `gorm:"column:id;primaryKey"`
	Email     string    `gorm:"column:email;uniqueIndex"`
	Password  string    `gorm:"column:password"`
	Token     string    `gorm:"column:token"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (u *User) TableName() string {
	return "users"
}
