package db

import (
	"os"
	"testing"

	"github.com/gobuffalo/packr"
	log "github.com/sirupsen/logrus"
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
		packr.NewBox("./migrations"),
	)
	if cErr != nil {
		log.Fatal(cErr)
	}

	os.Exit(m.Run())
}
