package main

import (
	"log"

	"github.com/alex-ant/ports/config"
	"github.com/alex-ant/ports/source"
)

func main() {
	// Read source file.
	_, srErr := source.NewReader(*config.SourceFile)
	if srErr != nil {
		log.Fatalf("failed to init source file reader: %v", srErr)
	}

}
