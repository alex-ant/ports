package main

import (
	"log"

	"github.com/alex-ant/ports/config"
	"github.com/alex-ant/ports/db"
)

func main() {
	dcClient, dcClientErr := db.New(
		*config.DBUser,
		*config.DBPass,
		*config.DBHost,
		*config.DBPort,
		*config.DBName,
		*config.DBTimeout,
	)
	if dcClientErr != nil {
		log.Fatalf("failed to establish DB connection: %v", dcClientErr)
	}

	defer dcClient.Close()
}
