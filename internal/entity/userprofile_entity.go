package entity

import "time"

type UserProfile struct {
	ID         uint              `gorm:"column:id;primaryKey"`
	UserId     uint              `gorm:"column:user_id"`
	Gender     UserProfileGender `gorm:"column:gender"`
	Name       string            `gorm:"column:name"`
	PictureURL string            `gorm:"column:picture_url"`
	CreatedAt  time.Time         `gorm:"column:created_at"`
	UpdatedAt  time.Time         `gorm:"column:updated_at"`
}

func (u *UserProfile) TableName() string {
	return "user_profiles"
}

type UserProfileGender string

const UserProfileGenderMale = "male"
const UserProfileGenderFemale = "female"
