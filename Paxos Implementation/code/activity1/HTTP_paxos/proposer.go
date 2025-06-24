package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Prepare struct {
	ProposalNumber int
}
type Promise struct {
	ProposalNumber int
	AcceptedValue  interface{}
}
type Accept struct {
	ProposalNumber int
	Value          interface{}
}
type Accepted struct {
	ProposalNumber int
	Value          interface{}
}

type Proposer struct {
	ProposalNumber int
	Value          interface{}
}

func (p *Proposer) sendRequest(url string, payload interface{}, response interface{}) error {
	data, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(response)
}

func (p *Proposer) Propose(value interface{}, acceptorURLs []string) interface{} {
	promises := 0
	majority := len(acceptorURLs) / 2 // Majority is more than half

	// Step 1: Send prepare requests and collect promises
	for _, url := range acceptorURLs {
		var promise Promise
		err := p.sendRequest(url+"/prepare", Prepare{ProposalNumber: p.ProposalNumber}, &promise)
		if promise.ProposalNumber == p.ProposalNumber {
			promises++
			fmt.Printf("Promise received from %s: %+v\n", url, promise)
		} else {
			fmt.Printf("Promise rejected by %s: Error: %v\n", url, err)
		}
	}

	// Step 2: Ensure majority of promises before moving to accept phase
	if promises > majority {
		fmt.Printf("Majority promises received (%d/%d).\n", promises, len(acceptorURLs))

		accepted := 0
		// Step 3: Send accept requests to all acceptors
		for _, url := range acceptorURLs {
			var ack Accepted
			err := p.sendRequest(url+"/accept", Accept{ProposalNumber: p.ProposalNumber, Value: value}, &ack)
			if err == nil && ack.ProposalNumber == p.ProposalNumber {
				accepted++
				fmt.Printf("Accept received from %s: %+v\n", url, ack)
			} else {
				fmt.Printf("Accept rejected by %s: Error: %v\n", url, err)
			}
		}

		// Step 4: If majority acceptances are reached, consensus is achieved
		if accepted > majority {
			fmt.Println("Consensus reached!")
			return value
		}
		fmt.Printf("Consensus failed after accept phase (%d/%d).\n", accepted, len(acceptorURLs))
	} else {
		fmt.Printf("Not enough promises (%d/%d) to proceed to accept phase.\n", promises, len(acceptorURLs))
	}

	fmt.Println("Consensus not reached.")
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: proposer <acceptor1> <acceptor2> ...")
		return
	}

	acceptorURLs := os.Args[1:]

	proposer := Proposer{ProposalNumber: 1, Value: "Distributed Systems"}
	value := proposer.Propose("Distributed Systems - HTTP communication", acceptorURLs)
	if value != nil {
		fmt.Printf("Consensus reached on value: %s\n", value)
	} else {
		fmt.Println("Consensus not reached")
	}
}
