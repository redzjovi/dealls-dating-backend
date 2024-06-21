package model

import "time"

type UpdateUserProfileRequest struct {
	Gender     string `json:"gender" validate:"required"`
	Name       string `json:"name" validate:"required"`
	PictureURL string `json:"picture_url" validate:"required"`
}

type UserProfileResponse struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	Gender     string    `json:"gender"`
	Name       string    `json:"name"`
	PictureURL string    `json:"picture_url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
