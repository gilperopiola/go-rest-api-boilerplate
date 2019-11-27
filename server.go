package main

import (
	"flag"
	"log"
	"os"

	"github.com/gilperopiola/go-rest-api-boilerplate/config"
)

var cfg config.MyConfig
var db MyDatabase
var rtr MyRouter

func main() {
	env := flag.String("env", "local", "local / dev / prod")
	flag.Parse()

	cfg.Setup(*env)
	db.Setup(cfg)
	defer db.Close()
	rtr.Setup(true)

	log.Println("server started")
	rtr.Run(":" + os.Getenv("PORT"))
}
