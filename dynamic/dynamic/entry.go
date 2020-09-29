package dynamic

type Entry interface {
	GetT() *Timestamp
	GetCmd() *Command
	GetStatus() int
	SetStatus(int)
}

type CmdEntry struct {
	Cmd    *Command
	T      *Timestamp
	Status int
}

func (e *CmdEntry) GetT() *Timestamp {
	return e.T
}

func (e *CmdEntry) GetCmd() *Command {
	return e.Cmd
}

func (e *CmdEntry) GetStatus() int {
	return e.Status
}

func (e *CmdEntry) SetStatus(s int) {
	e.Status = s
}
