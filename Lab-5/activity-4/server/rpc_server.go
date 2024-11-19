package main

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"sync"
)

// Args holds the arguments for arithmetic operations
type Args struct {
	A, B int
}

type Calculator struct {
	lastResult int
	mu         sync.Mutex
}

// GetLastResult retrieves the last result
func (c *Calculator) GetLastResult(args *Args, reply *int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	*reply = c.lastResult
	return nil
}

// Multiply multiplies two integers and returns the result
func (c *Calculator) Multiply(args *Args, reply *int) error {
	if args.A == 0 || args.B == 0 {
		return errors.New("multiplication by zero is not allowed")
	}
	*reply = args.A * args.B
	c.lastResult = *reply
	return nil
}

// Add adds two integers and returns the result
func (c *Calculator) Add(args *Args, reply *int) error {
	*reply = args.A + args.B
	c.lastResult = *reply
	return nil
}

// Subtract subtracts two integers and returns the result
func (c *Calculator) Subtract(args *Args, reply *int) error {
	*reply = args.A - args.B
	c.lastResult = *reply
	return nil
}

// Divide divides two integers and returns the result
func (c *Calculator) Divide(args *Args, reply *int) error {
	if args.B == 0 {
		return errors.New("division by zero!")
	}
	*reply = args.A / args.B
	c.lastResult = *reply
	return nil
}

func main() {
	// Register the Calculator service
	calc := new(Calculator)
	var serv rpc.Server
	serv.Register(calc)

	// Start listening for incoming RPC connections
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("Error starting RPC server:", err)
		return
	}

	fmt.Println("RPC server is listening on port 1234...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go serv.ServeConn(conn) // Handle each client in a goroutine
	}
}
