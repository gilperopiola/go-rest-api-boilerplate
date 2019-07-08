package main

import (
	"testing"

	"github.com/gilperopiola/go-rest-api-boilerplate/utils"
	"github.com/stretchr/testify/assert"
)

func createTestUserStruct(identifier int) *User {
	stringIdentifier := utils.ToString(identifier)

	return &User{
		Email:     "email" + stringIdentifier,
		Password:  hash("email"+stringIdentifier, "password"),
		FirstName: "firstName" + stringIdentifier,
		LastName:  "lastName" + stringIdentifier,
		Roles:     []Role{RoleUser},
	}
}

func TestLoginUser(t *testing.T) {
	cfg.Setup("test")
	db.Setup(cfg)
	defer db.Close()

	user := createTestUserStruct(1)
	user, _ = user.Create()

	user, err := user.Login()
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.Equal(t, "email1", user.Email)
	assert.Equal(t, "firstName1", user.FirstName)
	assert.Equal(t, "lastName1", user.LastName)
	assert.Equal(t, true, user.Enabled)
	assert.NotZero(t, user.DateCreated)
	assert.NotZero(t, user.Token)
}

func TestCreateUser(t *testing.T) {
	cfg.Setup("test")
	db.Setup(cfg)
	defer db.Close()

	user := createTestUserStruct(1)
	user, err := user.Create()

	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.Equal(t, "email1", user.Email)
	assert.Equal(t, "firstName1", user.FirstName)
	assert.Equal(t, "lastName1", user.LastName)
	assert.Equal(t, RoleUser, user.Roles[0])
	assert.Equal(t, true, user.Enabled)
	assert.NotZero(t, user.DateCreated)
}

func TestGetUser(t *testing.T) {
	cfg.Setup("test")
	db.Setup(cfg)
	defer db.Close()

	user := createTestUserStruct(1)
	user, _ = user.Create()

	user, err := user.Get()
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.Equal(t, "email1", user.Email)
	assert.Equal(t, "firstName1", user.FirstName)
	assert.Equal(t, "lastName1", user.LastName)
	assert.Equal(t, true, user.Enabled)
	assert.NotZero(t, user.DateCreated)
}

func TestSearchUsers(t *testing.T) {
	cfg.Setup("test")
	db.Setup(cfg)
	defer db.Close()

	user := createTestUserStruct(1)
	user2 := createTestUserStruct(2)
	user3 := createTestUserStruct(3)
	user4 := createTestUserStruct(4)
	user, _ = user.Create()
	user2, _ = user2.Create()
	user3, _ = user3.Create()
	user4, _ = user4.Create()

	params := &UserSearchParameters{
		FilterID:        "",
		FilterEmail:     "",
		FilterFirstName: "",
		FilterLastName:  "",

		SortField:     "id",
		SortDirection: "DESC",

		Limit:  3,
		Offset: 1,
	}

	users, err := user.Search(params)
	assert.NoError(t, err)

	assert.Equal(t, params.Limit, len(users))
	assert.NotZero(t, users[0].ID)
	assert.Equal(t, "email3", users[0].Email)
	assert.Equal(t, "firstName3", users[0].FirstName)
	assert.Equal(t, "lastName3", users[0].LastName)
	assert.True(t, users[0].Enabled)
	assert.NotZero(t, users[0].DateCreated)
}

func TestUpdateUser(t *testing.T) {
	cfg.Setup("test")
	db.Setup(cfg)
	defer db.Close()

	user := createTestUserStruct(1)
	user, _ = user.Create()

	user.Email = "email2"
	user.FirstName = "firstName2"
	user.Roles = []Role{RoleAdmin}

	user, err := user.Update()
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.Equal(t, "email2", user.Email)
	assert.Equal(t, "firstName2", user.FirstName)
	assert.Equal(t, RoleAdmin, user.Roles[0])

	user, err = user.ToggleEnabled()
	assert.NoError(t, err)
	assert.False(t, user.Enabled)

	user, err = user.ToggleEnabled()
	assert.NoError(t, err)
	assert.True(t, user.Enabled)
}
