package test

import (
	"dealls-dating/internal/entity"
	"dealls-dating/internal/model"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindUserProfileNotFound(t *testing.T) {
	ClearAll()
	TestAuthLogin(t)

	user := new(entity.User)
	err := db.Where("email = ?", "email@email.com").Take(user).Error
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodGet, "/api/auth/user/profile", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserProfileResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
	assert.NotNil(t, responseBody.Errors)
}

func TestFindUserProfile(t *testing.T) {
	ClearAll()
	TestUpdateUserProfile(t)

	user := new(entity.User)
	err := db.Where("email = ?", "email@email.com").Take(user).Error
	assert.Nil(t, err)

	userProfile := new(entity.UserProfile)
	err = db.Where("user_id = ?", user.ID).Take(userProfile).Error
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodGet, "/api/auth/user/profile", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserProfileResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, string(userProfile.Gender), responseBody.Data.Gender)
	assert.Equal(t, userProfile.Name, responseBody.Data.Name)
	assert.Equal(t, userProfile.PictureURL, responseBody.Data.PictureURL)
}

func TestUpdateUserProfileBadRequest(t *testing.T) {
	ClearAll()
	TestAuthLogin(t)

	user := new(entity.User)
	err := db.Where("email = ?", "email@email.com").Take(user).Error
	assert.Nil(t, err)

	requestBody := model.UpdateUserProfileRequest{}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPut, "/api/auth/user/profile", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserProfileResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.NotNil(t, responseBody.Errors)
}

func TestUpdateUserProfile(t *testing.T) {
	ClearAll()
	TestAuthLogin(t)

	user := new(entity.User)
	err := db.Where("email = ?", "email@email.com").Take(user).Error
	assert.Nil(t, err)

	requestBody := model.UpdateUserProfileRequest{
		Gender:     "male",
		Name:       "name",
		PictureURL: "picture_url",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPut, "/api/auth/user/profile", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNoContent, response.StatusCode)
}
