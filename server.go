package main

import (
	"flag"
	"log"
	"os"

	"github.com/gilperopiola/go-rest-api-boilerplate/config"
	"github.com/gilperopiola/go-rest-api-boilerplate/database"
)

var cfg config.MyConfig
var db database.MyDatabase
var rtr MyRouter

func main() {
	env := flag.String("env", "local", "local / dev / prod")
	flag.Parse()

	cfg.Setup(*env)
	db.Setup(cfg)
	defer db.Close()
	rtr.Setup()

	log.Println("server started")
	rtr.Run(":" + os.Getenv("PORT"))
}
