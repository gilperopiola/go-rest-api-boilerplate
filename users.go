package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gilperopiola/frutils"
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
	Login() (*User, error)

	Create() (*User, error)
	Get() (*User, error)
	Update() (*User, error)
	Search() ([]*User, error)
}

type UserTestingActions interface {
	GenerateTestRequest(token, method, url string) *httptest.ResponseRecorder
	GenerateJSONBody() string

	getRolesJSONBodyString() string
	GenerateSearchURLString() string
}

//UserActions
func (user *User) Login() (*User, error) {
	if err := db.DB.QueryRow(`SELECT id FROM users WHERE email = ? AND password = ?`, user.Email, user.Password).Scan(&user.ID); err != nil {
		return &User{}, err
	}

	user.Token = generateToken(*user)
	return user, nil
}

func (user *User) Create() (*User, error) {
	result, err := db.DB.Exec(`INSERT INTO users (email, password, firstName, lastName) VALUES (?, ?, ?, ?)`, user.Email, user.Password, user.FirstName, user.LastName)
	if err != nil {
		return &User{}, err
	}

	user.ID = frutils.GetID(result)

	user.createRoles()

	return user.Get()
}

func (user *User) Search(params *UserSearchParameters) ([]*User, error) {

	query := fmt.Sprintf(`SELECT id, email, firstName, lastName, enabled, dateCreated 
						  FROM users WHERE id LIKE ? AND email LIKE ? AND firstName LIKE ? AND lastName LIKE ?
						  ORDER BY %s LIMIT ? OFFSET ?`, frutils.GetOrderByString(params.SortField, params.SortDirection))

	rows, err := db.DB.Query(query, params.getQueryFormat()...)

	defer rows.Close()
	if err != nil {
		return []*User{}, err
	}

	users := []*User{}
	for rows.Next() {
		tempUser := &User{}

		if err = rows.Scan(&tempUser.ID, &tempUser.Email, &tempUser.FirstName, &tempUser.LastName, &tempUser.Enabled, &tempUser.DateCreated); err != nil {
			return []*User{}, err
		}

		if tempUser.Roles, err = tempUser.getRoles(); err != nil {
			return []*User{}, err
		}

		users = append(users, tempUser)
	}

	return users, nil
}

func (user *User) Get() (*User, error) {
	err := db.DB.QueryRow(`SELECT email, password, firstName, lastName, enabled, dateCreated FROM users WHERE id = ?`, user.ID).Scan(&user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Enabled, &user.DateCreated)
	if err != nil {
		return &User{}, err
	}

	if user.Roles, err = user.getRoles(); err != nil {
		return &User{}, err
	}

	return user, nil
}

func (user *User) Update() (*User, error) {
	if _, err := db.DB.Exec(`UPDATE users SET email = ?, firstName = ?, lastName = ? WHERE id = ?`, user.Email, user.FirstName, user.LastName, user.ID); err != nil {
		return &User{}, err
	}

	if err := user.updateRoles(); err != nil {
		return &User{}, err
	}

	return user.Get()
}

func (user *User) ToggleEnabled() (*User, error) {
	if _, err := db.DB.Exec(`UPDATE users SET enabled = NOT enabled WHERE id = ?`, user.ID); err != nil {
		return &User{}, err
	}

	return user.Get()
}

//UserTestingActions
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
		"enabled": ` + frutils.BoolToString(user.Enabled) + `
	}`
	return body
}

func (params *UserSearchParameters) generateSearchURLString() string {
	return fmt.Sprintf("?id=%s&email=%s&firstName=%s&lastName=%s&sortField=%s&sortDirection=%s&limit=%d&offset=%d", params.FilterID, params.FilterEmail, params.FilterFirstName, params.FilterLastName, params.SortField, params.SortDirection, params.Limit, params.Offset)
}
