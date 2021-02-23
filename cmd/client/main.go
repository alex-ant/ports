package main

import (
	"fmt"
	"log"

	"github.com/alex-ant/ports/config"
	"github.com/alex-ant/ports/port"
	"github.com/alex-ant/ports/source"
)

func main() {
	// Read source file.
	sr, srErr := source.NewReader(*config.SourceFile)
	if srErr != nil {
		log.Fatalf("failed to init source file reader: %v", srErr)
	}

	readErr := sr.Read(func(id string, pi port.Info) error {
		fmt.Println(id, pi)
		return nil
	})
	if readErr != nil {
		log.Fatalf("failed to read source file: %v", readErr)
	}
}
