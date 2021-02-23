package config

import (
	"flag"

	"github.com/alex-ant/envs"
	log "github.com/sirupsen/logrus"
)

var (
	APIPort = flag.Int("api-port", 8080, "API port")

	SourceFile = flag.String("source-file", "ports.json", "Ports source file")

	DBUser    = flag.String("db-user", "postgres", "DB user")
	DBPass    = flag.String("db-pass", "postgres", "DB pass")
	DBHost    = flag.String("db-host", "ports-test-postgres", "DB host")
	DBPort    = flag.Int("db-port", 5432, "DB port")
	DBName    = flag.String("db-name", "ports", "DB name")
	DBTimeout = flag.Int("db-timeout", 30, "DB connection timeout in seconds")

	GRPCPort       = flag.Int("grpc-port", 9000, "gPRC connection port")
	PortDomainHost = flag.String("port-domain-host", "port-domain", "gRPC port domain host")
)

// Parse parses the incomning flags and extracts the environment variables.
func init() {
	// Parse flags if not parsed already.
	if !flag.Parsed() {
		flag.Parse()
	}

	// Determine and read environment variables.
	flagsErr := envs.GetAllFlags()
	if flagsErr != nil {
		log.Fatal(flagsErr)
	}
}
