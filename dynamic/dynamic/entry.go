package dynamic

type Entry interface {
	GetT() *Timestamp
	GetCmd() *Command
	GetStatus() int
	SetStatus(int)
	SetStartDuration(int64)
	GetStartDuration() int64
}

type CmdEntry struct {
	Cmd    *Command
	T      *Timestamp
	Status int
	Duration int64
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

func (e *CmdEntry) SetStartDuration(t int64){
	e.Duration = t
}

func (e *CmdEntry) GetStartDuration() int64{
	return e.Duration
}