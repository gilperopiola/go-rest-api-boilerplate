package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginUser(t *testing.T) {
	cfg.Setup("_test")
	db.Setup(cfg)
	defer db.Close()

	user := &User{
		Email:     "email",
		Password:  "password",
		FirstName: "first_name",
		LastName:  "last_name",
		Roles:     []Role{RoleUser},
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

func TestCreateUser(t *testing.T) {
	cfg.Setup("_test")
	db.Setup(cfg)
	defer db.Close()

	user := &User{
		Email:     "email",
		Password:  "password",
		FirstName: "first_name",
		LastName:  "last_name",
		Roles:     []Role{RoleUser},
	}

	user, err := user.Create()
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.Equal(t, "email", user.Email)
	assert.Equal(t, "password", user.Password)
	assert.Equal(t, "first_name", user.FirstName)
	assert.Equal(t, "last_name", user.LastName)
	assert.Equal(t, RoleUser, user.Roles[0])
	assert.Equal(t, true, user.Enabled)
	assert.NotZero(t, user.DateCreated)
}

func TestGetUser(t *testing.T) {
	cfg.Setup("_test")
	db.Setup(cfg)
	defer db.Close()

	user := &User{
		Email:     "email",
		Password:  "password",
		FirstName: "first_name",
		LastName:  "last_name",
		Roles:     []Role{RoleUser},
	}
	user, _ = user.Create()

	user, err := user.Get()
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.Equal(t, "email", user.Email)
	assert.Equal(t, "password", user.Password)
	assert.Equal(t, "first_name", user.FirstName)
	assert.Equal(t, "last_name", user.LastName)
	assert.Equal(t, true, user.Enabled)
	assert.NotZero(t, user.DateCreated)
}

func TestGetAllUsers(t *testing.T) {
	cfg.Setup("_test")
	db.Setup(cfg)
	defer db.Close()

	user := &User{
		Email:     "email",
		Password:  "password",
		FirstName: "first_name",
		LastName:  "last_name",
		Roles:     []Role{RoleUser},
	}
	user, _ = user.Create()

	users, err := user.GetAll()
	assert.NoError(t, err)
	assert.NotZero(t, users[0].ID)
	assert.Equal(t, "email", users[0].Email)
	assert.Equal(t, "first_name", users[0].FirstName)
	assert.Equal(t, "last_name", users[0].LastName)
	assert.True(t, users[0].Enabled)
	assert.NotZero(t, users[0].DateCreated)
}

func TestUpdateUser(t *testing.T) {
	cfg.Setup("_test")
	db.Setup(cfg)
	defer db.Close()

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

	user, err := user.Update()
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.Equal(t, "email2", user.Email)
	assert.Equal(t, "first_name2", user.FirstName)
	assert.Equal(t, RoleAdmin, user.Roles[0])

	user, err = user.ToggleEnabled()
	assert.NoError(t, err)
	assert.False(t, user.Enabled)

	user, err = user.ToggleEnabled()
	assert.NoError(t, err)
	assert.True(t, user.Enabled)
}
