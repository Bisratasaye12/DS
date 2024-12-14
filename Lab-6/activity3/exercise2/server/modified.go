package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

// Args holds the arguments for arithmetic operations
type Args struct {
	A, B int
}

// Calculator is the type that will handle the RPC methods
type Calculator struct{}

// Divide is the method that performs division
func (c *Calculator) Divide(args *Args, reply *int) error {
	// Handle division by zero
	if args.B == 0 {
		return fmt.Errorf("cannot divide by zero")
	}
	*reply = args.A / args.B
	return nil
}

func main() {
	calculator := new(Calculator)

	// Register the Calculator type for RPC
	err := rpc.Register(calculator)
	if err != nil {
		log.Fatal("Error registering Calculator:", err)
	}

	// Create a listener on TCP port 1234
	listener, err := net.Listen("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Error starting listener:", err)
	}
	defer listener.Close()

	fmt.Println("Server is listening on localhost:1234...")
	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Error accepting connection:", err)
		}

		// Handle RPC requests on the connection
		go rpc.ServeConn(conn)
	}
}
