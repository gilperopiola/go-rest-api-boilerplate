package main

import (
	"time"

	"github.com/gilperopiola/lyfe-companyon-backend/utils"
)

type User struct {
	ID          int
	Email       string
	Password    string
	FirstName   string
	LastName    string
	Enabled     bool
	DateCreated time.Time

	Token string
}

type UserActions interface {
	Create() (*User, error)
	GetByID() (*User, error)

	Login() (*User, error)
}

func (user *User) Create() (*User, error) {
	result, err := db.DB.Exec(`INSERT INTO users (email, password, first_name, last_name) VALUES (?, ?, ?, ?)`, user.Email, user.Password, user.FirstName, user.LastName)
	if err != nil {
		return &User{}, err
	}

	user.ID = utils.StripLastInsertID(result.LastInsertId())

	user, err = user.GetByID()
	if err != nil {
		return &User{}, err
	}

	return user, nil
}

func (user *User) GetByID() (*User, error) {
	if err := db.DB.QueryRow(`SELECT email, password, first_name, last_name, enabled, date_created FROM users WHERE id = ?`, user.ID).Scan(
		&user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Enabled, &user.DateCreated); err != nil {
		return &User{}, err
	}

	return user, nil
}

func (user *User) Login() (*User, error) {
	if err := db.DB.QueryRow(`SELECT id FROM users WHERE email = ? AND password = ?`, user.Email, user.Password).Scan(&user.ID); err != nil {
		return &User{}, err
	}

	user.Token = generateToken(*user)
	return user, nil
}
