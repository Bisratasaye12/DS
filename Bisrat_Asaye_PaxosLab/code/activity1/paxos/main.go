package main

import (
	"fmt"
)

func main() {
	// Simulate 5 acceptors (2k+1 with k=2)
	acceptors := []*Acceptor{
		&Acceptor{},
		&Acceptor{},
		&Acceptor{},
		&Acceptor{},
		&Acceptor{},
	}

	proposer := Proposer{ProposalNumber: 1, Value: "Distributed Systems"}

	// Simulate consensus without failure
	fmt.Println("Testing without failure:")
	value := proposer.Propose("Distributed Systems", acceptors)
	if value != nil {
		fmt.Printf("Consensus reached on value: %s\n", value)
	} else {
		fmt.Println("Consensus not reached")
	}

	// Simulate failures
	fmt.Println("\nTesting with failures (2 acceptors down):")
	activeAcceptors := acceptors[:3] // remove 2 acceptors
	value = proposer.Propose("Fault Tolerant Systems", activeAcceptors)
	if value != nil {
		fmt.Printf("Consensus reached on value: %s\n", value)
	} else {
		fmt.Println("Consensus not reached")
		// }

	}

}
