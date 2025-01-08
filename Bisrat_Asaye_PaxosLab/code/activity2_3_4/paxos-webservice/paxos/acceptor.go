// package paxos

// import (
// 	"sync"
// )

// type Acceptor struct {
// 	mu             sync.Mutex
// 	promisedNumber int
// 	acceptedNumber int
// 	acceptedValue  interface{}
// }

// func (a *Acceptor) HandlePrepare(p Prepare) Promise {
// 	a.mu.Lock()
// 	defer a.mu.Unlock()
// 	if p.ProposalNumber > a.promisedNumber {
// 		a.promisedNumber = p.ProposalNumber
// 		return Promise{ProposalNumber: p.ProposalNumber,
// 			AcceptedValue: a.acceptedValue}
// 	}
// 	return Promise{}
// }
// func (a *Acceptor) HandleAccept(ac Accept) Accepted {
// 	a.mu.Lock()
// 	defer a.mu.Unlock()
// 	if ac.ProposalNumber >= a.promisedNumber {
// 		a.promisedNumber = ac.ProposalNumber
// 		a.acceptedNumber = ac.ProposalNumber
// 		a.acceptedValue = ac.Value
// 		return Accepted{ProposalNumber: ac.ProposalNumber, Value: ac.Value}
// 	}
// 	return Accepted{}
// }

package paxos

import (
	"context"
	"fmt"
	"sync"
)

// Acceptor handles the prepare and accept phases of the Paxos protocol
type Acceptor struct {
	mu             sync.Mutex
	promisedNumber int
	acceptedNumber int
	acceptedValue  interface{}
}

// HandlePrepareWithContext handles the prepare phase, considering context for timeouts
func (a *Acceptor) HandlePrepareWithContext(ctx context.Context, p Prepare) Promise {
	select {
	case <-ctx.Done():
		fmt.Println("Timeout during HandlePrepare")
		return Promise{}
	default:
		a.mu.Lock()
		defer a.mu.Unlock()
		if p.ProposalNumber > a.promisedNumber {
			a.promisedNumber = p.ProposalNumber
			return Promise{ProposalNumber: p.ProposalNumber, AcceptedValue: a.acceptedValue}
		}
		return Promise{}
	}
}

// HandleAcceptWithContext handles the accept phase, considering context for timeouts
func (a *Acceptor) HandleAcceptWithContext(ctx context.Context, ac Accept) Accepted {
	select {
	case <-ctx.Done():
		fmt.Println("Timeout during HandleAccept")
		return Accepted{}
	default:
		a.mu.Lock()
		defer a.mu.Unlock()
		if ac.ProposalNumber >= a.promisedNumber {
			a.promisedNumber = ac.ProposalNumber
			a.acceptedNumber = ac.ProposalNumber
			a.acceptedValue = ac.Value
			return Accepted{ProposalNumber: ac.ProposalNumber, Value: ac.Value}
		}
		return Accepted{}
	}
}
