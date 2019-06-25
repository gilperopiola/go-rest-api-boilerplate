package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUserController(t *testing.T) {
	cfg.Setup("_test")
	db.Setup(cfg)
	defer db.Close()
	rtr.Setup()

	token := generateTestingToken()

	user := User{Email: "newUser@email.com", Password: "password"}

	response := user.GenerateTestRequest(token, "POST", "")
	json.Unmarshal(response.Body.Bytes(), &user)
	assert.Equal(t, http.StatusOK, response.Code)

	assert.Equal(t, "newUser@email.com", user.Email)
	assert.Equal(t, true, user.Enabled)
}
