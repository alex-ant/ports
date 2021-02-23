package config

import (
	"flag"
	"log"

	"github.com/alex-ant/envs"
)

var (
	SourceFile = flag.String("source-file", "ports.json", "Ports source file")
)

// Parse parses the incomning flags and extracts the environment variables.
func Parse() {
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
