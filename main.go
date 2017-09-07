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

// TODO: Fuck it. Move everything to main.go...
// This likely won't be a HUGE service. If it becomes that, just extract things out.

func init() {
	// Register all codecs.
	// Later on we could probably just include the ones we need?
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
	server, err := streaming.NewStreamingServer(db)
	if err != nil {
		log.Panic(err)
	}

	// Initialize our Environment.
	// Setup our Database tables for testing
	env := &Env{db, server}
	models.InitTables(env.db)

	// Start the streaming server.
	go env.stream.Start()

	// HTTP Handler functions
	http.HandleFunc("/", env.stream.HandleHTTP)

	// Start the HTTP Server and listen
	http.ListenAndServe(":8089", nil)
}
