package model

type SwipeLikeResponse struct {
	Match bool `json:"match"`
}
type SwipeUserResponse struct {
	UserId     uint   `json:"user_id"`
	Gender     string `json:"gender"`
	Name       string `json:"name"`
	PictureURL string `json:"picture_url"`
}
