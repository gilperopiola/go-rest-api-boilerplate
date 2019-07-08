package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUserController(t *testing.T) {
	cfg.Setup("test")
	db.Setup(cfg)
	defer db.Close()
	rtr.Setup()
	token := generateTestingToken(RoleAdmin)

	user := createTestUserStruct(1)

	response := user.GenerateTestRequest(token, "POST", "")
	json.Unmarshal(response.Body.Bytes(), &user)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "email1", user.Email)
	assert.Equal(t, true, user.Enabled)
	assert.Equal(t, RoleUser, user.Roles[0])
}

func TestSearchUserController(t *testing.T) {
	cfg.Setup("test")
	db.Setup(cfg)
	defer db.Close()
	rtr.Setup()
	token := generateTestingToken(RoleAdmin)

	user := createTestUserStruct(1)
	user2 := createTestUserStruct(2)
	user, _ = user.Create()
	user2, _ = user2.Create()

	params := &UserSearchParameters{
		FilterID:        "",
		FilterEmail:     "1",
		FilterFirstName: "1",
		FilterLastName:  "1",

		SortField:     "id",
		SortDirection: "DESC",

		Limit:  100,
		Offset: 0,
	}

	users := make([]*User, 0)
	response := user.GenerateTestRequest(token, "GET", params.generateSearchURLString())
	json.Unmarshal(response.Body.Bytes(), &users)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, 1, len(users))
	assert.Equal(t, "email1", users[0].Email)
	assert.Equal(t, "firstName1", users[0].FirstName)
	assert.Equal(t, RoleUser, users[0].Roles[0])
}

func TestGetUserController(t *testing.T) {
	cfg.Setup("test")
	db.Setup(cfg)
	defer db.Close()
	rtr.Setup()
	token := generateTestingToken(RoleAdmin)

	user := createTestUserStruct(1)
	user, _ = user.Create()

	response := user.GenerateTestRequest(token, "GET", "/"+strconv.Itoa(user.ID))
	json.Unmarshal(response.Body.Bytes(), &user)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "email1", user.Email)
}

func TestUpdateUserController(t *testing.T) {
	cfg.Setup("test")
	db.Setup(cfg)
	defer db.Close()
	rtr.Setup()
	token := generateTestingToken(RoleAdmin)

	user := createTestUserStruct(1)
	user, _ = user.Create()

	user.Email = "email2"
	user.FirstName = "firstName2"
	user.Roles = []Role{RoleAdmin}

	response := user.GenerateTestRequest(token, "PUT", "/"+strconv.Itoa(user.ID))
	json.Unmarshal(response.Body.Bytes(), &user)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "email2", user.Email)
	assert.Equal(t, "firstName2", user.FirstName)
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
