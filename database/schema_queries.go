package database

var createUsersTableQuery = `
CREATE TABLE IF NOT EXISTS users (
	id int UNIQUE NOT NULL AUTO_INCREMENT,
	email varchar(255) UNIQUE NOT NULL,
	password varchar(255) NOT NULL,
	first_name varchar(255),
	last_name varchar(255),
	enabled tinyint DEFAULT '1',
	date_created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
)
`
