package main

import (
	"net/http"

	"github.com/gilperopiola/frutils"
	"github.com/gin-gonic/gin"
)

//Signup takes {email, password, repeatPassword}. It creates a user and returns it.
func Signup(c *gin.Context) {
	var auth Auth
	if bindErr := c.BindJSON(&auth); bindErr != nil {
		c.JSON(http.StatusBadRequest, bindErr.Error())
		return
	}

	if auth.Email == "" || auth.Password == "" || auth.RepeatPassword == "" {
		c.JSON(http.StatusBadRequest, "all fields required")
		return
	}

	if auth.Password != auth.RepeatPassword {
		c.JSON(http.StatusBadRequest, "passwords don't match")
		return
	}

	hashedPassword := hash(auth.Email, auth.Password)

	user := &User{
		Email:    auth.Email,
		Password: hashedPassword,
	}

	user, err := user.Create()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

//Login takes {email, password}, checks if the user exists and returns it.
func Login(c *gin.Context) {
	var auth Auth
	if bindErr := c.BindJSON(&auth); bindErr != nil {
		c.JSON(http.StatusBadRequest, bindErr.Error())
		return
	}

	if auth.Email == "" || auth.Password == "" {
		c.JSON(http.StatusBadRequest, "all fields required")
		return
	}

	hashedPassword := hash(auth.Email, auth.Password)

	user := &User{
		Email:    auth.Email,
		Password: hashedPassword,
	}

	user, err := user.Login()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context) {
	var user *User
	if bindErr := c.BindJSON(&user); bindErr != nil {
		c.JSON(http.StatusBadRequest, bindErr.Error())
		return
	}

	if user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, "email and password required")
		return
	}

	user.Password = hash(user.Email, user.Password)

	if len(user.Roles) == 0 {
		user.Roles = []Role{RoleUser}
	}

	user, err := user.Create()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func GetUser(c *gin.Context) {
	user := &User{ID: frutils.ToInt(c.Param("id_user"))}

	user, err := user.Get()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	var user *User
	if bindErr := c.BindJSON(&user); bindErr != nil {
		c.JSON(http.StatusBadRequest, bindErr.Error())
		return
	}

	if user.Email == "" {
		c.JSON(http.StatusBadRequest, "email required")
		return
	}

	user.ID = frutils.ToInt(c.Param("id_user"))

	user, err := user.Update()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func ToggleUserEnabled(c *gin.Context) {
	user := &User{ID: frutils.ToInt(c.Param("id_user"))}

	user, err := user.ToggleEnabled()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func SearchUsers(c *gin.Context) {
	user := &User{}

	params := &UserSearchParameters{}
	params = params.Fill(c)

	users, err := user.Search(params)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, users)
}
