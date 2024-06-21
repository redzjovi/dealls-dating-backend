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

func TestAuthSignUpBadRequest(t *testing.T) {
	ClearAll()
	requestBody := model.AuthSignUpRequest{}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/sign-up", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestAuthSignUpConflict(t *testing.T) {
	ClearAll()
	TestAuthSignUp(t)
	requestBody := model.AuthSignUpRequest{
		Email:    "email@email.com",
		Password: "password",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/sign-up", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusConflict, response.StatusCode)
}

func TestAuthSignUp(t *testing.T) {
	ClearAll()
	requestBody := model.AuthSignUpRequest{
		Email:    "email@email.com",
		Password: "password",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/sign-up", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNoContent, response.StatusCode)
}

func TestAuthLoginUnauthorized(t *testing.T) {
	ClearAll()
	TestAuthSignUp(t)

	requestBody := model.AuthLoginRequest{
		Email:    "email2@email.com",
		Password: "password",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.AuthLoginResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.NotNil(t, responseBody.Errors)
}

func TestAuthLoginUnauthorizedPassword(t *testing.T) {
	ClearAll()
	TestAuthSignUp(t)

	requestBody := model.AuthLoginRequest{
		Email:    "email@email.com",
		Password: "password2",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.AuthLoginResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.NotNil(t, responseBody.Errors)
}

func TestAuthLogin(t *testing.T) {
	ClearAll()
	TestAuthSignUp(t)

	requestBody := model.AuthLoginRequest{
		Email:    "email@email.com",
		Password: "password",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.AuthLoginResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data.Token)

	user := new(entity.User)
	err = db.Where("email = ?", requestBody.Email).Take(user).Error
	assert.Nil(t, err)
	assert.Equal(t, user.Token, responseBody.Data.Token)
}

func TestAuthUserLogoutUnauthorized(t *testing.T) {
	ClearAll()
	TestAuthLogin(t)

	request := httptest.NewRequest(http.MethodDelete, "/api/auth/user", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "wrong")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[bool])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.NotNil(t, responseBody.Errors)
}

func TestAuthUserLogout(t *testing.T) {
	ClearAll()
	TestAuthLogin(t)

	user := new(entity.User)
	err := db.Where("email = ?", "email@email.com").Take(user).Error
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodDelete, "/api/auth/user", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNoContent, response.StatusCode)
}
