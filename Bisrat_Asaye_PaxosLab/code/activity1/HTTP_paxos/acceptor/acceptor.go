package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
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

type Acceptor struct {
	mu             sync.Mutex
	promisedNumber int
	acceptedNumber int
	acceptedValue  interface{}
}

func (a *Acceptor) HandlePrepare(w http.ResponseWriter, r *http.Request) {
	var p Prepare
	json.NewDecoder(r.Body).Decode(&p)
	a.mu.Lock()
	defer a.mu.Unlock()
	if p.ProposalNumber > a.promisedNumber {
		a.promisedNumber = p.ProposalNumber
		response := Promise{ProposalNumber: p.ProposalNumber, AcceptedValue: a.acceptedValue}
		json.NewEncoder(w).Encode(response)
		return
	}
	json.NewEncoder(w).Encode(Promise{})
}

func (a *Acceptor) HandleAccept(w http.ResponseWriter, r *http.Request) {
	var ac Accept
	json.NewDecoder(r.Body).Decode(&ac)
	a.mu.Lock()
	defer a.mu.Unlock()
	if ac.ProposalNumber >= a.promisedNumber {
		a.promisedNumber = ac.ProposalNumber
		a.acceptedNumber = ac.ProposalNumber
		a.acceptedValue = ac.Value
		response := Accepted{ProposalNumber: ac.ProposalNumber, Value: ac.Value}
		json.NewEncoder(w).Encode(response)
		return
	}

	json.NewEncoder(w).Encode(Accepted{})
}

func main() {
	// Get the port from command-line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: acceptor <port>")
		return
	}
	port := os.Args[1]

	// Create the acceptor instance
	acceptor := &Acceptor{}
	http.HandleFunc("/prepare", acceptor.HandlePrepare)
	http.HandleFunc("/accept", acceptor.HandleAccept)

	fmt.Printf("Acceptor running on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Error starting acceptor: %v\n", err)
	}
}
