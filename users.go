package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	"github.com/gilperopiola/go-rest-api-boilerplate/utils"
)

type User struct {
	ID          int
	Email       string
	Password    string
	FirstName   string
	LastName    string
	Roles       []Role
	Enabled     bool
	DateCreated time.Time

	Token string
}

type UserActions interface {
	//External
	Login() (*User, error)
	Create() (*User, error)
	Get() (*User, error)
	GetAll() ([]*User, error)
	Update() (*User, error)
	Delete() (*User, error)

	GenerateTestRequest(token, method, url string) *httptest.ResponseRecorder
	GenerateJSONBody() string

	//Internal
	createRoles() error
	getRoles() ([]Role, error)
	updateRoles() ([]Role, error)
	deleteRoles() error

	getRolesJSONBodyString() string
}

type Role int

const (
	RoleUser  Role = 0
	RoleAdmin Role = 1
)

//External Functions
func (user *User) Login() (*User, error) {
	if err := db.DB.QueryRow(`SELECT id FROM users WHERE email = ? AND password = ?`, user.Email, user.Password).Scan(&user.ID); err != nil {
		return &User{}, err
	}

	user.Token = generateToken(*user)
	return user, nil
}

func (user *User) Create() (*User, error) {
	result, err := db.DB.Exec(`INSERT INTO users (email, password, first_name, last_name) VALUES (?, ?, ?, ?)`, user.Email, user.Password, user.FirstName, user.LastName)
	if err != nil {
		return &User{}, err
	}

	user.ID = utils.StripLastInsertID(result.LastInsertId())

	user.createRoles()

	return user.Get()
}

func (user *User) Get() (*User, error) {
	err := db.DB.QueryRow(`SELECT email, password, first_name, last_name, enabled, date_created FROM users WHERE id = ?`, user.ID).
		Scan(&user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Enabled, &user.DateCreated)
	if err != nil {
		return &User{}, err
	}

	user.Roles, err = user.getRoles()
	if err != nil {
		return &User{}, err
	}

	return user, nil
}

func (user *User) GetAll() ([]*User, error) {
	rows, err := db.DB.Query(`SELECT id, email, first_name, last_name, enabled, date_created FROM users`)
	defer rows.Close()
	if err != nil {
		return []*User{}, err
	}

	users := []*User{}
	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Enabled, &user.DateCreated)
		if err != nil {
			return []*User{}, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (user *User) Update() (*User, error) {
	_, err := db.DB.Exec(`UPDATE users SET email = ?, first_name = ?, last_name = ? WHERE id = ?`, user.Email, user.FirstName, user.LastName, user.ID)
	if err != nil {
		return &User{}, err
	}

	err = user.updateRoles()
	if err != nil {
		return &User{}, err
	}

	return user, nil
}

func (user *User) ToggleEnabled() (*User, error) {
	_, err := db.DB.Exec(`UPDATE users SET enabled = NOT enabled WHERE id = ?`, user.ID)
	if err != nil {
		return &User{}, err
	}

	return user.Get()
}

func (user *User) GenerateTestRequest(token, method, url string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	body := user.GetJSONBody()
	req, _ := http.NewRequest(method, "/User"+url, bytes.NewReader([]byte(body)))
	req.Header.Set("Authorization", token)
	rtr.ServeHTTP(w, req)
	return w
}

func (user *User) GetJSONBody() string {
	body := `{
		"email": "` + user.Email + `",
		"password": "` + user.Password + `",
		"firstName": "` + user.FirstName + `",
		"lastName": "` + user.LastName + `",
		"roles": [` + user.getRolesJSONBodyString() + `],
		"enabled": ` + utils.BoolToString(user.Enabled) + `
	}`
	return body
}

//Internal Functions

func (user *User) createRoles() error {
	for _, role := range user.Roles {
		_, err := db.DB.Exec(`INSERT INTO users_roles (id_user, id_role) VALUES (?, ?)`, user.ID, role)
		if err != nil {
			return err
		}
	}

	return nil
}

func (user *User) getRoles() ([]Role, error) {
	rows, err := db.DB.Query(`SELECT id_role FROM users_roles WHERE id_user = ?`, user.ID)
	defer rows.Close()
	if err != nil {
		return []Role{}, err
	}

	roles := []Role{}
	for rows.Next() {
		var role Role
		err = rows.Scan(&role)
		if err != nil {
			return []Role{}, err
		}

		roles = append(roles, role)
	}

	return roles, nil
}

func (user *User) updateRoles() error {
	err := user.deleteRoles()
	if err != nil {
		return err
	}

	for _, role := range user.Roles {
		_, err = db.DB.Exec(`INSERT INTO users_roles (id_user, id_role) VALUES (?, ?)`, user.ID, role)
		if err != nil {
			return err
		}
	}

	return nil
}

func (user *User) deleteRoles() error {
	_, err := db.DB.Exec(`DELETE FROM users_roles WHERE id_user = ?`, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (user *User) getRolesJSONBodyString() string {
	rolesString := ""

	for i, role := range user.Roles {
		rolesString += strconv.Itoa(int(role))

		if i+1 < len(user.Roles) {
			rolesString += ", "
		}
	}

	return rolesString
}
