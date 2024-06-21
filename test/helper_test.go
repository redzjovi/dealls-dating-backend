package test

import (
	"dealls-dating/internal/entity"
)

func ClearAll() {
	ClearUserProfiles()
	ClearUsers()
}

func ClearUserProfiles() {
	err := db.Where("id is not null").Delete(&entity.UserProfile{}).Error
	if err != nil {
		log.Fatalf("Failed clear user profile data : %+v", err)
	}
}

func ClearUsers() {
	err := db.Where("id is not null").Delete(&entity.User{}).Error
	if err != nil {
		log.Fatalf("Failed clear user data : %+v", err)
	}
}
