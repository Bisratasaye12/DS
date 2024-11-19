package main

import (
	"fmt"
	"log"
	"net/rpc"
)

// Args holds the arguments for multiplication
type Args struct {
	A, B int
}

func main() {
	// Connect to the RPC server
	client, err := rpc.Dial("tcp", "localhost:1235")
	if err != nil {
		log.Fatal("Error connecting to RPC server:", err)
	}
	// Prepare the arguments and call the Multiply method
	
	args := Args{A: 3, B: 5}
	var reply int
	err = client.Call("Calculator.Multiply", &args, &reply)
	if err != nil {
		log.Fatal("Error calling RPC:", err)
	}
	fmt.Printf("Result of %d * %d = %d\n", args.A, args.B, reply)

	err = client.Call("Calculator.Add", &args, &reply)
	if err != nil {
		log.Fatal("error calling rpc, add method")
	}
	fmt.Printf("Result of %d + %d = %d\n", args.A, args.B, reply)

	err = client.Call("Calculator.Subtract", &args, &reply)
	if err != nil {
		log.Fatal("error calling rpc, subtract")
	}
	fmt.Printf("Result of %d + %d = %d\n", args.A, args.B, reply)

	err = client.Call("Calculator.Divide", &args, &reply)
	if err != nil {
		log.Fatal("error calling rpc, divide method")
	}
	fmt.Printf("Result of %d / %d = %d\n", args.A, args.B, reply)

	
}


