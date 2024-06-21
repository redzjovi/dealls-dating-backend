package test

import (
	"dealls-dating/internal/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindSwipePreconditionFailed(t *testing.T) {
	ClearAll()
	user := CreateUser(t)

	request := httptest.NewRequest(http.MethodGet, "/api/auth/user/swipe", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.SwipeUserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusPreconditionFailed, response.StatusCode)
	assert.NotNil(t, responseBody.Errors)
}

func TestFindSwipeNotFound(t *testing.T) {
	ClearAll()
	CreateUserProfile(t, 0)
	user := GetFirstUser(t)

	request := httptest.NewRequest(http.MethodGet, "/api/auth/user/swipe", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestFindSwipe(t *testing.T) {
	ClearAll()
	CreateUserProfiles(t, 10)
	user := GetFirstUser(t)

	request := httptest.NewRequest(http.MethodGet, "/api/auth/user/swipe", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestLikeSwipeTooManyRequests(t *testing.T) {
	ClearAll()
	userProfiles := CreateUserProfiles(t, 12)
	userProfiles = userProfiles[0:11]
	user12 := GetLastUser(t)

	for i, userProfile := range userProfiles {
		request := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/auth/user/swipe/%v", userProfile.UserId), nil)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")
		request.Header.Set("Authorization", user12.Token)

		response, err := app.Test(request)
		assert.Nil(t, err)

		if i+1 > 10 {
			assert.Equal(t, http.StatusTooManyRequests, response.StatusCode)
		} else {
			assert.Equal(t, http.StatusOK, response.StatusCode)
		}
	}
}

func TestLikeSwipePremium(t *testing.T) {
	ClearAll()
	userProfiles := CreateUserProfiles(t, 11)
	CreateUserPremium(t)
	user12 := GetLastUser(t)

	for _, userProfile := range userProfiles {
		request := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/auth/user/swipe/%v", userProfile.UserId), nil)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")
		request.Header.Set("Authorization", user12.Token)

		response, err := app.Test(request)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, response.StatusCode)
	}
}

func TestLikeSwipeConflict(t *testing.T) {
	ClearAll()
	CreateUserProfiles(t, 2)
	user1 := GetFirstUser(t)
	user2 := GetLastUser(t)

	for i := 1; i <= 2; i++ {
		request := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/auth/user/swipe/%v", user1.ID), nil)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")
		request.Header.Set("Authorization", user2.Token)

		response, err := app.Test(request)
		assert.Nil(t, err)

		if i == 1 {
			assert.Equal(t, http.StatusOK, response.StatusCode)
		} else if i > 1 {
			assert.Equal(t, http.StatusConflict, response.StatusCode)
		}
	}
}

func TestLikeSwipe(t *testing.T) {
	ClearAll()
	CreateUserProfiles(t, 2)
	user1 := GetFirstUser(t)
	user2 := GetLastUser(t)

	request := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/auth/user/swipe/%v", user1.ID), nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user2.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.SwipeLikeResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, false, responseBody.Data.Match)
}

func TestLikeSwipeAndMatch(t *testing.T) {
	ClearAll()
	CreateUserProfiles(t, 2)
	user1 := GetFirstUser(t)
	user2 := GetLastUser(t)

	// 1st
	request := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/auth/user/swipe/%v", user2.ID), nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user1.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.SwipeLikeResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, false, responseBody.Data.Match)

	// 2nd
	request = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/auth/user/swipe/%v", user1.ID), nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user2.Token)

	response, err = app.Test(request)
	assert.Nil(t, err)

	bytes, err = io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody = new(model.WebResponse[model.SwipeLikeResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, true, responseBody.Data.Match)
}

func TestDisLikeSwipeTooManyRequests(t *testing.T) {
	ClearAll()
	userProfiles := CreateUserProfiles(t, 12)
	userProfiles = userProfiles[0:11]
	user := GetLastUser(t)

	for i, userProfile := range userProfiles {
		request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/auth/user/swipe/%v", userProfile.UserId), nil)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")
		request.Header.Set("Authorization", user.Token)

		response, err := app.Test(request)
		assert.Nil(t, err)

		if i+1 > 10 {
			assert.Equal(t, http.StatusTooManyRequests, response.StatusCode)
		} else {
			assert.Equal(t, http.StatusNoContent, response.StatusCode)
		}
	}
}

func TestDisLikeSwipePremium(t *testing.T) {
	ClearAll()
	userProfiles := CreateUserProfiles(t, 11)
	CreateUserPremium(t)
	user := GetLastUser(t)

	for _, userProfile := range userProfiles {
		request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/auth/user/swipe/%v", userProfile.UserId), nil)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")
		request.Header.Set("Authorization", user.Token)

		response, err := app.Test(request)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusNoContent, response.StatusCode)
	}
}

func TestDislikeSwipeConflict(t *testing.T) {
	ClearAll()
	CreateUserProfiles(t, 2)
	user1 := GetFirstUser(t)
	user2 := GetLastUser(t)

	for i := 1; i <= 2; i++ {
		request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/auth/user/swipe/%v", user1.ID), nil)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")
		request.Header.Set("Authorization", user2.Token)

		response, err := app.Test(request)
		assert.Nil(t, err)

		if i == 1 {
			assert.Equal(t, http.StatusNoContent, response.StatusCode)
		} else if i > 1 {
			assert.Equal(t, http.StatusConflict, response.StatusCode)
		}
	}
}

func TestDislikeSwipe(t *testing.T) {
	ClearAll()
	CreateUserProfiles(t, 2)
	user1 := GetFirstUser(t)
	user2 := GetLastUser(t)

	request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/auth/user/swipe/%v", user1.ID), nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user2.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNoContent, response.StatusCode)
}
