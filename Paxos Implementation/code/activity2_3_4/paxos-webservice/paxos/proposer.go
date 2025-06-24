// package paxos

// type Proposer struct {
// 	ProposalNumber int
// 	Value          interface{}
// }

// func (p *Proposer) Propose(value interface{}, acceptors []*Acceptor) interface{} {
// 	promises := 0
// 	for _, acceptor := range acceptors {
// 		promise := acceptor.HandlePrepare(Prepare{ProposalNumber: p.ProposalNumber})
// 		if promise.ProposalNumber == p.ProposalNumber {
// 			promises++
// 		}
// 	}
// 	if promises > len(acceptors)/2 {
// 		accepted := 0
// 		for _, acceptor := range acceptors {
// 			ack := acceptor.HandleAccept(Accept{ProposalNumber: p.ProposalNumber, Value: value})
// 			if ack.ProposalNumber == p.ProposalNumber {
// 				accepted++
// 			}
// 		}
// 		if accepted > len(acceptors)/2 {
// 			return value
// 		}
// 	}
// 	return nil
// }

package paxos

import (
	"context"
	"log"
	"time"
)

type Proposer struct {
	ProposalNumber int
	Value          interface{}
}

// Propose attempts to propose a value with retries and context-based timeouts
func (p *Proposer) Propose(value interface{}, acceptors []*Acceptor) interface{} {
	maxRetries := 3
	for retry := 1; retry <= maxRetries; retry++ {
		log.Printf("Attempt %d to propose value %v with proposal number %d", retry, value, p.ProposalNumber)

		// Attempt proposal
		result := p.attemptProposal(value, acceptors)
		if result != nil {
			return result
		}

		// Log and retry
		log.Printf("Proposal failed. Retrying... (Attempt %d)", retry)
		time.Sleep(2 * time.Second) // Delay before retrying
	}

	log.Printf("Proposal %d failed after %d attempts", p.ProposalNumber, maxRetries)
	return nil
}

// attemptProposal makes a single attempt to propose a value to the acceptors
func (p *Proposer) attemptProposal(value interface{}, acceptors []*Acceptor) interface{} {
	// Create a context with timeout for communication
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	promises := 0
	for _, acceptor := range acceptors {
		// Send the prepare request to each acceptor and check if it's successful
		promise := acceptor.HandlePrepareWithContext(ctx, Prepare{ProposalNumber: p.ProposalNumber})
		if promise.ProposalNumber == p.ProposalNumber {
			promises++
		} else {
			log.Printf("Failed communication attempt: Failed to receive valid promise from acceptor")
		}
	}

	if promises > len(acceptors)/2 {
		accepted := 0
		for _, acceptor := range acceptors {
			// Send the accept request to each acceptor and check if it's successful
			ack := acceptor.HandleAcceptWithContext(ctx, Accept{ProposalNumber: p.ProposalNumber, Value: value})
			if ack.ProposalNumber == p.ProposalNumber {
				accepted++
			} else {
				log.Printf("Failed communication attempt: Failed to receive valid acceptance from acceptor")
			}
		}
		if accepted > len(acceptors)/2 {
			log.Printf("Proposal %d accepted successfully with value: %v", p.ProposalNumber, value)
			return value
		}
	}
	return nil
}
