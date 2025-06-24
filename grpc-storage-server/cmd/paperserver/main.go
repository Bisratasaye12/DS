package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	pb "paper-storage-server/paper-storage-server/paperpb"
	"paper-storage-server/server"
	"sync"
	"syscall"

	"github.com/streadway/amqp"
	"google.golang.org/grpc"
)

func main() {
	// Establish a connection to RabbitMQ
	rabbitConn, err := connectRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitConn.Close()

	// Start the gRPC server
	fmt.Printf("\nPaper Storage gRPC server is running...\n\n")
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Initialize the server with necessary fields
	s := grpc.NewServer()

	server := &server.Server{
		RabbitConn:   rabbitConn,
		ContentStore: make(map[int32][]byte),            // Initialize content store
		DetailStore:  make(map[int32]map[string]string), // Initialize detail store
		ContentMutex: sync.Mutex{},                      // Initialize mutex for content store
		DetailMutex:  sync.Mutex{},                      // Initialize mutex for detail store
		ID:           0,
	}

	pb.RegisterPaperStorageServiceServer(s, server)

	// Graceful shutdown handling
	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for an interrupt signal (Ctrl+C)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	sigReceived := <-signalChan
	fmt.Printf("\nReceived signal: %s. Shutting down...\n", sigReceived)

	// Graceful shutdown after receiving a termination signal
	s.GracefulStop()
}

// connectRabbitMQ establishes a connection to RabbitMQ and returns the connection object
func connectRabbitMQ() (*amqp.Connection, error) {
	// Establish the connection
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	// Set up a signal listener for graceful shutdown
	go func() {
		// Capture OS interrupt signals
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

		<-sigs // Wait for a signal
		log.Println("Shutting down RabbitMQ connection...")
		if err := conn.Close(); err != nil {
			log.Printf("Failed to close RabbitMQ connection: %v", err)
		} else {
			log.Println("RabbitMQ connection closed successfully.")
		}
		os.Exit(0) // Exit the application
	}()

	log.Println("Successfully connected to RabbitMQ")
	return conn, nil
}
