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
	args := Args{A: 10, B: 0} // Division by zero to trigger an error
	var reply int
	// Call RPC method with a timeout
	call := client.Go("Calculator.Add", &args, &reply, nil)
	//call := client.Go("Calculator.Subtract", &args, &reply, nil)
	//call := client.Go("Calculator.Multiply", &args, &reply, nil)
	//call := client.Go("Calculator.Divide", &args, &reply, nil)
	
	select {
	case <-call.Done:
		if call.Error != nil {
			log.Println("RPC error:", call.Error)
		} else {
			fmt.Printf("Result: %d\n", reply)
		}
	case <-time.After(2 * time.Second):
		log.Println("RPC call timed out")
	}
}
