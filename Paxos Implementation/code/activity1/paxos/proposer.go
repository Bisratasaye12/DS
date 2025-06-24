package main


type Proposer struct {
	ProposalNumber int
	Value          interface{}
}

func (p *Proposer) Propose(value interface{}, acceptors []*Acceptor) interface{} {
	promises := 0
	for _, acceptor := range acceptors {
		promise := acceptor.HandlePrepare(Prepare{ProposalNumber: p.ProposalNumber})
		if promise.ProposalNumber == p.ProposalNumber {
			promises++
		}
	}
	if promises > len(acceptors)/2 {
		accepted := 0
		for _, acceptor := range acceptors {
			ack := acceptor.HandleAccept(Accept{ProposalNumber: p.ProposalNumber, Value: value})
			if ack.ProposalNumber == p.ProposalNumber {
				accepted++
			}
		}
		if accepted > len(acceptors)/2 {
			return value
		}
	}
	return nil
}
