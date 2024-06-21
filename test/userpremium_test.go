package test

import (
	"dealls-dating/internal/model"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListUserPremium(t *testing.T) {
	ClearAll()
	user := CreateUser(t)

	request := httptest.NewRequest(http.MethodGet, "/api/auth/user/premium", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[[]model.UserPremiumResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestTrialUserPremiumConflict(t *testing.T) {
	ClearAll()
	CreateUserPremium(t)
	user := GetFirstUser(t)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/user/premium", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusConflict, response.StatusCode)
}

func TestTrialUserPremium(t *testing.T) {
	ClearAll()
	user := CreateUser(t)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/user/premium", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNoContent, response.StatusCode)
}
