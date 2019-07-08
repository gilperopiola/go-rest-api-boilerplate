package database

import (
	"database/sql"
	"log"
	"strings"

	"github.com/gilperopiola/go-rest-api-boilerplate/config"
	_ "github.com/go-sql-driver/mysql"
)

type DatabaseActions interface {
	Setup()
	Purge()
	Migrate()

	LoadTestingData()
	BeautifyError(error) string
}

type MyDatabase struct {
	*sql.DB
}

func (db *MyDatabase) Setup(cfg config.MyConfig) {
	var err error
	db.DB, err = sql.Open(
		cfg.DATABASE.TYPE, cfg.DATABASE.USERNAME+":"+cfg.DATABASE.PASSWORD+"@tcp("+cfg.DATABASE.HOSTNAME+":"+
			cfg.DATABASE.PORT+")/"+cfg.DATABASE.SCHEMA+"?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	err = db.DB.Ping()
	if err != nil {
		log.Fatalf("error pinging database: %v", err)
	}

	if cfg.DATABASE.CREATE_SCHEMA {
		db.CreateSchema()
	}

	if cfg.DATABASE.PURGE {
		db.Purge()
	}

	if cfg.DATABASE.CREATE_ADMIN {
		db.CreateAdmin()
	}
}

func (db *MyDatabase) CreateSchema() {
	if _, err := db.DB.Exec(createUsersTableQuery); err != nil {
		log.Println(err.Error())
	}

	if _, err := db.DB.Exec(createUsersRolesTableQuery); err != nil {
		log.Println(err.Error())
	}
}

func (db *MyDatabase) Purge() {
	db.DB.Exec("DELETE FROM users")
	db.DB.Exec("DELETE FROM users_roles")
}

func (db *MyDatabase) CreateAdmin() {
	result, err := db.DB.Exec("INSERT INTO users (email, password, firstName, lastName) VALUES ('ferra.main@gmail.com', 'MD1pQbDskXpNndv7zuLsZJt24RY=', 'Franco', 'Ferraguti')")
	if err != nil {
		log.Println(err)
		return
	}

	id, _ := result.LastInsertId()

	if _, err := db.DB.Exec("INSERT INTO users_roles (idUser, idRole) VALUES (?, 1)", id); err != nil {
		log.Println(err)
	}
}

func (db *MyDatabase) BeautifyError(err error) string {
	s := err.Error()

	if strings.Contains(s, "Duplicate entry") {
		duplicateField := strings.Split(s, "'")[3]
		return duplicateField + " already in use"
	}

	return s
}
