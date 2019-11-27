package main

var createUsersTableQuery = `
CREATE TABLE IF NOT EXISTS users (
	id int UNIQUE NOT NULL AUTO_INCREMENT,
	email varchar(255) UNIQUE NOT NULL,
	password varchar(255) NOT NULL,
	firstName varchar(255),
	lastName varchar(255),
	enabled tinyint DEFAULT '1',
	dateCreated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
)
`

var createUsersRolesTableQuery = `
CREATE TABLE IF NOT EXISTS users_roles (
	id int UNIQUE NOT NULL AUTO_INCREMENT,
	idUser int NOT NULL,
	idRole int NOT NULL
)
`
