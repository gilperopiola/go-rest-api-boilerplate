package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUserController(t *testing.T) {
	cfg.Setup("_test")
	db.Setup(cfg)
	defer db.Close()
	rtr.Setup()
	token := generateTestingToken(RoleAdmin)

	user := User{Email: "newUser@email.com", Password: "password"}

	response := user.GenerateTestRequest(token, "POST", "")
	json.Unmarshal(response.Body.Bytes(), &user)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "newUser@email.com", user.Email)
	assert.Equal(t, true, user.Enabled)
	assert.Equal(t, RoleUser, user.Roles[0])
}

func TestSearchUserController(t *testing.T) {
	cfg.Setup("_test")
	db.Setup(cfg)
	defer db.Close()
	rtr.Setup()
	token := generateTestingToken(RoleAdmin)

	user := &User{
		Email:     "email",
		Password:  "password",
		FirstName: "first_name",
		LastName:  "last_name",
		Roles:     []Role{RoleUser},
	}
	user, _ = user.Create()

	users := make([]*User, 0)
	response := user.GenerateTestRequest(token, "GET", "")
	json.Unmarshal(response.Body.Bytes(), &users)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "email", users[0].Email)
	assert.Equal(t, "first_name", users[0].FirstName)
	assert.Equal(t, RoleUser, users[0].Roles[0])
}

func TestGetUserController(t *testing.T) {
	cfg.Setup("_test")
	db.Setup(cfg)
	defer db.Close()
	rtr.Setup()
	token := generateTestingToken(RoleAdmin)

	user := &User{
		Email:     "email",
		Password:  "password",
		FirstName: "first_name",
		LastName:  "last_name",
		Roles:     []Role{RoleUser},
	}
	user, _ = user.Create()

	response := user.GenerateTestRequest(token, "GET", "/"+strconv.Itoa(user.ID))
	json.Unmarshal(response.Body.Bytes(), &user)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "email", user.Email)
}

func TestUpdateUserController(t *testing.T) {
	cfg.Setup("_test")
	db.Setup(cfg)
	defer db.Close()
	rtr.Setup()
	token := generateTestingToken(RoleAdmin)

	user := &User{
		Email:     "email",
		Password:  "password",
		FirstName: "first_name",
		LastName:  "last_name",
		Roles:     []Role{RoleUser},
	}
	user, _ = user.Create()

	user.Email = "email2"
	user.FirstName = "first_name2"
	user.Roles = []Role{RoleAdmin}

	response := user.GenerateTestRequest(token, "PUT", "/"+strconv.Itoa(user.ID))
	json.Unmarshal(response.Body.Bytes(), &user)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "email2", user.Email)
	assert.Equal(t, "first_name2", user.FirstName)
	assert.Equal(t, RoleAdmin, user.Roles[0])

	response = user.GenerateTestRequest(token, "PUT", "/"+strconv.Itoa(user.ID)+"/Enabled")
	json.Unmarshal(response.Body.Bytes(), &user)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.False(t, user.Enabled)

	response = user.GenerateTestRequest(token, "PUT", "/"+strconv.Itoa(user.ID)+"/Enabled")
	json.Unmarshal(response.Body.Bytes(), &user)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.True(t, user.Enabled)
}
