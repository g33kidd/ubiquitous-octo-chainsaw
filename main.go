package main

import (
	"./models"
	"./streaming"

	"log"
	"net/http"

	"github.com/nareix/joy4/format"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func init() {
	format.RegisterAll()
}

// Env is our application environment
type Env struct {
	db     *gorm.DB
	stream *streaming.Server
}

func main() {

	// Open the database connection.
	db, err := models.NewDB("sqlite3", "./database/development.db")
	if err != nil {
		log.Panic(err)
	}

	// Setup the streaming server
	server, err := streaming.NewStreamingServer()
	if err != nil {
		log.Panic(err)
	}

	env := &Env{db, server}
	models.InitTables(env.db)

	go env.stream.Start()
	http.HandleFunc("/", env.stream.HandleHTTP)
	http.ListenAndServe(":8089", nil)
}
