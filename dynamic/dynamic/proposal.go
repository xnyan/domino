package dynamic

type PaxosProposal struct {
	Cmd  *Command
	RetC chan bool
}

type FpProposal struct {
	Cmd  *Command
	FpT  *Timestamp
	RetC chan bool
}
