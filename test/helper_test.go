package test

import (
	"dealls-dating/internal/entity"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
)

func ClearAll() {
	ClearMatches()
	ClearSwipes()
	ClearUserPremiums()
	ClearUserProfiles()
	ClearUsers()
}

func ClearMatches() {
	err := db.Where("id is not null").Delete(&entity.Match{}).Error
	if err != nil {
		log.Fatalf("Failed clear match data : %+v", err)
	}
}

func ClearSwipes() {
	err := db.Where("id is not null").Delete(&entity.Swipe{}).Error
	if err != nil {
		log.Fatalf("Failed clear swipe data : %+v", err)
	}
}

func ClearUserPremiums() {
	err := db.Where("id is not null").Delete(&entity.UserPremium{}).Error
	if err != nil {
		log.Fatalf("Failed clear user premium data : %+v", err)
	}
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

func CreateUser(t *testing.T) (res *entity.User) {
	user := entity.User{
		Email:    faker.Email(),
		Password: "password",
		Token:    faker.UUIDDigit(),
	}
	err := db.Create(&user).Error
	assert.Nil(t, err)
	return &user
}

func CreateUserPremium(t *testing.T) (res *entity.UserPremium) {
	user := CreateUser(t)
	userPremium := entity.UserPremium{
		UserId:  user.ID,
		StartAt: time.Now(),
		EndAt:   time.Now().Add(7 * 24 * time.Hour),
	}
	err := db.Create(&userPremium).Error
	assert.Nil(t, err)
	return &userPremium
}

func CreateUserProfile(t *testing.T, i int) *entity.UserProfile {
	user := entity.User{
		Email:    faker.Email(),
		Password: "password",
		Token:    faker.UUIDDigit(),
	}
	err := db.Create(&user).Error
	assert.Nil(t, err)
	gender := entity.UserProfileGenderMale
	if i%2 == 0 {
		gender = entity.UserProfileGenderFemale
	}
	userProfile := entity.UserProfile{
		UserId:     user.ID,
		Gender:     entity.UserProfileGender(gender),
		Name:       faker.Name(),
		PictureURL: faker.URL(),
	}
	err = db.Create(&userProfile).Error
	assert.Nil(t, err)
	return &userProfile
}

func CreateUserProfiles(t *testing.T, total int) (res []entity.UserProfile) {
	for i := 0; i < total; i++ {
		userProfile := CreateUserProfile(t, i)
		res = append(res, *userProfile)
	}
	return res
}

func GetFirstUser(t *testing.T) *entity.User {
	user := new(entity.User)
	err := db.First(user).Error
	assert.Nil(t, err)
	return user
}

func GetLastUser(t *testing.T) *entity.User {
	user := new(entity.User)
	err := db.Last(user).Error
	assert.Nil(t, err)
	return user
}
