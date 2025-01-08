package paxos

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
