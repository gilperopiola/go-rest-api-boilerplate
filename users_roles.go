package main

import (
	"strconv"
)

type Role int

const (
	RoleUser  Role = 0
	RoleAdmin Role = 1
)

func (user *User) createRoles() error {
	for _, role := range user.Roles {
		if _, err := db.DB.Exec(`INSERT INTO users_roles (idUser, idRole) VALUES (?, ?)`, user.ID, role); err != nil {
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

		if err = rows.Scan(&role); err != nil {
			return []Role{}, err
		}

		roles = append(roles, role)
	}

	return roles, nil
}

func (user *User) updateRoles() error {
	if err := user.deleteRoles(); err != nil {
		return err
	}

	for _, role := range user.Roles {
		if _, err := db.DB.Exec(`INSERT INTO users_roles (idUser, idRole) VALUES (?, ?)`, user.ID, role); err != nil {
			return err
		}
	}

	return nil
}

func (user *User) deleteRoles() error {
	if _, err := db.DB.Exec(`DELETE FROM users_roles WHERE idUser = ?`, user.ID); err != nil {
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
