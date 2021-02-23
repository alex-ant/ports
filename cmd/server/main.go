package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/alex-ant/ports/config"
	"github.com/alex-ant/ports/db"
	"github.com/alex-ant/ports/ports"
	"github.com/gobuffalo/packr"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	// Establish DB connection.
	dbClient, dbClientErr := db.New(
		*config.DBUser,
		*config.DBPass,
		*config.DBHost,
		*config.DBPort,
		*config.DBName,
		*config.DBTimeout,
		packr.NewBox("../../db/migrations"),
	)
	if dbClientErr != nil {
		log.Fatalf("failed to establish DB connection: %v", dbClientErr)
	}

	// Listen to gRPCs.
	lis, lisErr := net.Listen("tcp", fmt.Sprintf(":%d", *config.GRPCPort))
	if lisErr != nil {
		log.Fatalf("failed to listen on :%d: %v", *config.GRPCPort, lisErr)
	}

	grpcServer := grpc.NewServer()

	ps := ports.Server{}

	ports.RegisterPortServiceServer(grpcServer, &ps)

	go func() {
		serveErr := grpcServer.Serve(lis)
		if serveErr != nil {
			log.Fatalf("failed to serve gRPC: %v", serveErr)
		}
	}()

	// Shut down on SIGINT and SIGTERM.
	shutdown := func() {
		log.Println("Shutting down gracefully..")

		// Close database connection.
		dbClient.Close()

		// Stop gRPC server.
		grpcServer.Stop()

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
