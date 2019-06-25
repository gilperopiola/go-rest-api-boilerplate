package main

import (
	"net/http"

	"github.com/gilperopiola/go-rest-api-boilerplate/utils"

	"github.com/gin-gonic/gin"
)

//SignUp takes {email, password, repeatPassword}. It creates a user and returns it.
func SignUp(c *gin.Context) {
	var auth Auth
	c.BindJSON(&auth)

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
		c.JSON(http.StatusBadRequest, db.BeautifyError(err))
		return
	}

	c.JSON(http.StatusOK, user)
}

//Login takes {username, password}, checks if the user exists and returns it
func Login(c *gin.Context) {
	var auth Auth
	c.BindJSON(&auth)

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
		c.JSON(http.StatusBadRequest, db.BeautifyError(err))
		return
	}

	c.JSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context) {
	var user *User
	c.BindJSON(&user)

	if user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, "email and password required")
		return
	}

	user.Password = hash(user.Email, user.Password)

	user, err := user.Create()
	if err != nil {
		c.JSON(http.StatusBadRequest, db.BeautifyError(err))
		return
	}

	c.JSON(http.StatusOK, user)
}

func GetUser(c *gin.Context) {
	user := &User{ID: utils.ToInt(c.Param("id_user"))}

	user, err := user.Get()
	if err != nil {
		c.JSON(http.StatusBadRequest, db.BeautifyError(err))
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	var user *User
	c.BindJSON(&user)

	if user.Email == "" {
		c.JSON(http.StatusBadRequest, "email required")
		return
	}

	user.ID = utils.ToInt(c.Param("id_user"))

	user, err := user.Update()
	if err != nil {
		c.JSON(http.StatusBadRequest, db.BeautifyError(err))
		return
	}

	c.JSON(http.StatusOK, user)
}

func ToggleUserEnabled(c *gin.Context) {
	user := &User{ID: utils.ToInt(c.Param("id_user"))}

	user, err := user.ToggleEnabled()
	if err != nil {
		c.JSON(http.StatusBadRequest, db.BeautifyError(err))
		return
	}

	c.JSON(http.StatusOK, user)
}
