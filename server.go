package main

import (
	"log"
	"os"

	"github.com/gilperopiola/lyfe-companyon-backend/config"
	"github.com/gilperopiola/lyfe-companyon-backend/database"
)

var cfg config.MyConfig
var db database.MyDatabase
var rtr MyRouter

func main() {
	cfg.Setup("")
	db.Setup(cfg)
	defer db.Close()
	rtr.Setup()

	log.Println("server started")
	rtr.Run(":" + os.Getenv("PORT"))
}
