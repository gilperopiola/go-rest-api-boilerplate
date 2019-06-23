package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	cfg.Setup("_test")
	db.Setup(cfg)
	defer db.Close()

	user := &User{
		Email:     "email",
		Password:  "password",
		FirstName: "first_name",
		LastName:  "last_name",
	}

	user, err := user.Create()

	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.Equal(t, "email", user.Email)
	assert.Equal(t, "password", user.Password)
	assert.Equal(t, "first_name", user.FirstName)
	assert.Equal(t, "last_name", user.LastName)
	assert.Equal(t, true, user.Enabled)
	assert.NotZero(t, user.DateCreated)
}

func TestGetUserByID(t *testing.T) {
	cfg.Setup("_test")
	db.Setup(cfg)
	defer db.Close()

	user := &User{
		Email:     "email",
		Password:  "password",
		FirstName: "first_name",
		LastName:  "last_name",
	}

	user, _ = user.Create()
	user, err := user.GetByID()

	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.Equal(t, "email", user.Email)
	assert.Equal(t, "password", user.Password)
	assert.Equal(t, "first_name", user.FirstName)
	assert.Equal(t, "last_name", user.LastName)
	assert.Equal(t, true, user.Enabled)
	assert.NotZero(t, user.DateCreated)
}

func TestLoginUser(t *testing.T) {
	cfg.Setup("_test")
	db.Setup(cfg)
	defer db.Close()

	user := &User{
		Email:     "email",
		Password:  "password",
		FirstName: "first_name",
		LastName:  "last_name",
	}

	user, _ = user.Create()
	user, err := user.Login()

	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.Equal(t, "email", user.Email)
	assert.Equal(t, "password", user.Password)
	assert.Equal(t, "first_name", user.FirstName)
	assert.Equal(t, "last_name", user.LastName)
	assert.Equal(t, true, user.Enabled)
	assert.NotZero(t, user.DateCreated)
	assert.NotZero(t, user.Token)
}
