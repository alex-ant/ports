package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alex-ant/ports/api"
	"github.com/alex-ant/ports/config"
	"github.com/alex-ant/ports/ports"
	"github.com/alex-ant/ports/source"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	// Establish client gRPC connection.
	grpcConn, grpcConnErr := grpc.Dial(fmt.Sprintf("%s:%d", *config.PortDomainHost, *config.GRPCPort), grpc.WithInsecure())
	if grpcConnErr != nil {
		log.Fatalf("failed to establish client gRPC connection: %v", grpcConnErr)
	}

	c := ports.NewPortServiceClient(grpcConn)

	// Ping port domain, wait for it to become alive.
	for {
		_, pingErr := c.Ping(context.Background(), new(ports.Empty))
		if pingErr != nil {
			log.Println("gRPC port domain server ping failed, retrying in 1 second")
			time.Sleep(time.Second)
		} else {
			log.Println("gRPC connection established")
			break
		}
	}

	// Read source file.
	sr, srErr := source.NewReader(*config.SourceFile)
	if srErr != nil {
		log.Fatalf("failed to init source file reader: %v", srErr)
	}

	log.Println("Reading source JSON")

	var cout int
	readErr := sr.Read(func(pi *ports.PortInfo) error {
		_, storeErr := c.StorePortInfo(context.Background(), pi)
		cout++
		return storeErr
	})
	if readErr != nil {
		log.Fatalf("failed to read source file: %v", readErr)
	}

	log.Printf("Finished reading source JSON, %d records processed", cout)

	// Initialize API HTTP server.
	apiServer, apiServerErr := api.New(
		*config.APIPort,
		func() ([]*ports.PortInfo, error) {
			fetchedPortInfo, fetchedPortInfoErr := c.FetchPortInfo(context.Background(), new(ports.Empty))
			if fetchedPortInfoErr != nil {
				return nil, fmt.Errorf("failed to fetch port info: %v", fetchedPortInfoErr)
			}
			return fetchedPortInfo.Ports, nil
		},
	)
	if apiServerErr != nil {
		log.Fatal(apiServerErr)
	}

	// Start API HTTP server.
	apiStartErr := apiServer.Start()
	if apiStartErr != nil {
		log.Fatal(apiStartErr)
	}

	log.Printf("Started API server on port %d", *config.APIPort)

	// Shut down on SIGINT and SIGTERM.
	shutdown := func() {
		log.Println("Shutting down gracefully..")

		// Stop API server.
		apiServer.Stop()

		// Close gRPC client connection.
		grpcConn.Close()

		log.Println("terminating process")
		os.Exit(0)
	}

	go func() {
		intChan := make(chan os.Signal)
		signal.Notify(intChan, syscall.SIGINT, syscall.SIGTERM)
		<-intChan
		go shutdown()

		// Another signal will force process termination.
		signal.Notify(intChan, syscall.SIGINT, syscall.SIGTERM)
		<-intChan
		os.Exit(0)
	}()

	log.Println("Successfully started")

	// Keep the process running.
	select {}
}
