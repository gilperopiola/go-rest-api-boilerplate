package main

import (
	"bytes"
	"fmt"
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

//UserActions are external actions that the controllers can call
type UserActions interface {
	Login() (*User, error)

	Create() (*User, error)
	Search() ([]*User, error)
	Get() (*User, error)
	Update() (*User, error)
	Delete() (*User, error)
}

//UserInternals are the actual functions that work on the database
type UserInternals interface {
	createRoles() error
	getRoles() ([]Role, error)
	updateRoles() ([]Role, error)
	deleteRoles() error
}

//UserTestingActions are functions that aid in testing
type UserTestingActions interface {
	GenerateTestRequest(token, method, url string) *httptest.ResponseRecorder
	GenerateJSONBody() string
	getRolesJSONBodyString() string
	generateSearchURLString() string
}

type UserSearchParameters struct {
	FilterID        string
	FilterEmail     string
	FilterFirstName string
	FilterLastName  string

	SortField     string
	SortDirection string

	Limit  int
	Offset int
}

type Role int

const (
	RoleUser  Role = 0
	RoleAdmin Role = 1
)

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

	user.ID = utils.StripLastInsertID(result.LastInsertId())

	user.createRoles()

	return user.Get()
}

func (user *User) Search(params *UserSearchParameters) ([]*User, error) {
	orderByString := "id ASC"
	if params.SortField != "" && params.SortDirection != "" {
		orderByString = params.SortField + " " + params.SortDirection
	}

	query := fmt.Sprintf(`SELECT id, email, firstName, lastName, enabled, dateCreated 
						  FROM users WHERE id LIKE ? AND email LIKE ? AND firstName LIKE ? AND lastName LIKE ?
						  ORDER BY %s LIMIT ? OFFSET ?`, orderByString)

	rows, err := db.DB.Query(query, "%"+params.FilterID+"%", "%"+params.FilterEmail+"%", "%"+params.FilterFirstName+"%",
		"%"+params.FilterLastName+"%", params.Limit, params.Offset)

	defer rows.Close()
	if err != nil {
		return []*User{}, err
	}

	users := []*User{}
	for rows.Next() {
		tempUser := &User{}
		err = rows.Scan(&tempUser.ID, &tempUser.Email, &tempUser.FirstName, &tempUser.LastName, &tempUser.Enabled, &tempUser.DateCreated)
		if err != nil {
			return []*User{}, err
		}

		tempUser.Roles, err = tempUser.getRoles()
		if err != nil {
			return []*User{}, err
		}
		users = append(users, tempUser)
	}

	return users, nil
}

func (user *User) Get() (*User, error) {
	err := db.DB.QueryRow(`SELECT email, password, firstName, lastName, enabled, dateCreated FROM users WHERE id = ?`, user.ID).
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

func (user *User) Update() (*User, error) {
	_, err := db.DB.Exec(`UPDATE users SET email = ?, firstName = ?, lastName = ? WHERE id = ?`, user.Email, user.FirstName, user.LastName, user.ID)
	if err != nil {
		return &User{}, err
	}

	err = user.updateRoles()
	if err != nil {
		return &User{}, err
	}

	return user.Get()
}

func (user *User) ToggleEnabled() (*User, error) {
	_, err := db.DB.Exec(`UPDATE users SET enabled = NOT enabled WHERE id = ?`, user.ID)
	if err != nil {
		return &User{}, err
	}

	return user.Get()
}

//UserInternals
func (user *User) createRoles() error {
	for _, role := range user.Roles {
		_, err := db.DB.Exec(`INSERT INTO users_roles (idUser, idRole) VALUES (?, ?)`, user.ID, role)
		if err != nil {
			return err
		}
	}

	return nil
}

func (user *User) getRoles() ([]Role, error) {
	rows, err := db.DB.Query(`SELECT idRole FROM users_roles WHERE idUser = ?`, user.ID)
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
		_, err = db.DB.Exec(`INSERT INTO users_roles (idUser, idRole) VALUES (?, ?)`, user.ID, role)
		if err != nil {
			return err
		}
	}

	return nil
}

func (user *User) deleteRoles() error {
	_, err := db.DB.Exec(`DELETE FROM users_roles WHERE idUser = ?`, user.ID)
	if err != nil {
		return err
	}

	return nil
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
		"enabled": ` + utils.BoolToString(user.Enabled) + `
	}`
	return body
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

func (params *UserSearchParameters) generateSearchURLString() string {
	return fmt.Sprintf("?id=%s&email=%s&firstName=%s&lastName=%s&sortField=%s&sortDirection=%s&limit=%d&offset=%d",
		params.FilterID, params.FilterEmail, params.FilterFirstName, params.FilterLastName, params.SortField, params.SortDirection, params.Limit, params.Offset)
}
