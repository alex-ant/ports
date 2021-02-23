package db

import (
	"log"
	"os"
	"testing"
)

var testClient *Client

func TestMain(m *testing.M) {
	var cErr error
	testClient, cErr = New(
		"postgres",
		"postgres",
		"ports-test-postgres",
		5432,
		"ports",
		30,
	)
	if cErr != nil {
		log.Fatal(cErr)
	}

	os.Exit(m.Run())
}
