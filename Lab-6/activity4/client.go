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

	// Perform the Divide operation
	args := Args{A: 10, B: 2} // Replace with your operation arguments
	var reply int

	call := client.Go("Calculator.Divide", &args, &reply, nil)
	select {
	case <-call.Done:
		if call.Error != nil {
			log.Println("RPC error:", call.Error)
		} else {
			fmt.Printf("Divide Result: %d\n", reply)
		}
	case <-time.After(2 * time.Second): // Timeout of 2 seconds
		log.Println("RPC call timed out")
	}

	// Call GetLastResult to retrieve the last result
	var lastResult int
	err = client.Call("Calculator.GetLastResult", &args, &lastResult)
	if err != nil {
		log.Println("Error calling GetLastResult:", err)
	} else {
		fmt.Printf("Last Result from Server: %d\n", lastResult)
	}
}
