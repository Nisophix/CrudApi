package main

import (
	"log"

	app "github.com/Nisophix/crud_api/internal"
	"github.com/Nisophix/crud_api/internal/config"
)

func main() {
	config, err := config.ReadConfig("./config.toml")
	if err != nil {
		log.Fatal(err)
	}
	app := &app.App{}
	app.Initialize(config)
	app.Start(":8080")
}

// GetAll
// curl http://127.0.0.1:8080/events

// Create
// curl -X POST -H "Content-Type: application/json" -d '{"ID": 1, "Name": "Nicholas"}' http://127.0.0.1:8080/event

// Read
// curl http://127.0.0.1:8080/event/1

// Update
// curl -X PUT -H "Content-Type: application/json" -d '{"Name": "Joe"}' http://127.0.0.1:8080/event/1

// Delete
// curl -X DELETE http://127.0.0.1:8080/delete/1
