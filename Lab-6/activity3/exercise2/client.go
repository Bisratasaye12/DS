package main

import (
	"fmt"
	"log"
	"net/rpc"
	"time"
)

// Args holds the arguments for arithmetic operations
type Args struct {
	A, B int
}

func main() {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Error connecting to RPC server:", err)
	}
	defer client.Close()

	args := Args{A: 10, B: 2} // Replace with your operation arguments
	var reply int

	// Call RPC method with a timeout
	call := client.Go("Calculator.Divide", &args, &reply, nil)
	select {
	case <-call.Done:
		if call.Error != nil {
			log.Println("RPC error:", call.Error)
		} else {
			fmt.Printf("Result: %d\n", reply)
		}
	case <-time.After(2 * time.Second): // Timeout of 2 seconds
		log.Println("RPC call timed out")
	}
}
