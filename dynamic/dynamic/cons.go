package dynamic

// Paxos consenus instance
type PaxosCons struct {
	retC    chan bool
	acceptN int
}

func NewPaxosCons(c chan bool) *PaxosCons {
	return &PaxosCons{
		retC:    c,
		acceptN: 0,
	}
}

func (info *PaxosCons) VoteAccept() int {
	info.acceptN++
	return info.acceptN
}

func (info *PaxosCons) GetRetC() chan bool {
	return info.retC
}

// Fast Paxos consensus instance
type FpCons struct {
	t                   int64               // timestamp that identified this consensus instance
	n                   int                 // total number of accepted commands
	isDone              bool                // If a command is chosen
	cmd                 *Command            // the chosen command (nil for no-op)
	cmdMap              map[string]*Command // cmdId --> command
	cmdVoteMap          map[string]int      // cmdId --> accept N
	acceptCmdConsRetMap map[string]bool     // cmdId --> whether consensus result is set
}

func NewFpCons(t int64) *FpCons {
	return &FpCons{
		t:                   t,
		n:                   0,
		isDone:              false,
		cmd:                 nil,
		cmdMap:              make(map[string]*Command),
		cmdVoteMap:          make(map[string]int),
		acceptCmdConsRetMap: make(map[string]bool),
	}
}

func (cons *FpCons) GetT() int64 {
	return cons.t
}

func (cons *FpCons) Accept(cmd *Command) int {
	if cmd == nil {
		logger.Fatalf("Cannot accept a nil command")
	}

	cons.n++

	if _, ok := cons.acceptCmdConsRetMap[cmd.Id]; !ok {
		cons.acceptCmdConsRetMap[cmd.Id] = false
		cons.cmdMap[cmd.Id] = cmd
		cons.cmdVoteMap[cmd.Id] = 0
	}

	cons.cmdVoteMap[cmd.Id] += 1
	return cons.cmdVoteMap[cmd.Id]
}

func (cons *FpCons) ChooseCmd(cmd *Command) {
	if cons.isDone {
		logger.Fatalf("Cannot choose command twice. cmdId = %s", cmd.Id)
	}

	cons.isDone = true
	cons.cmd = cmd
}

// Returns the number of accepted commands
func (cons *FpCons) GetAcceptN() int {
	return cons.n
}

func (cons *FpCons) IsCmdChosen() bool {
	return cons.isDone
}

func (cons *FpCons) SelectCmd(majority int) *Command {
	for cmdId, cmd := range cons.cmdMap {
		if n, ok := cons.cmdVoteMap[cmdId]; ok {
			if n >= majority {
				return cmd
			}
		} else {
			logger.Fatalf("Cannot find vote info for cmdId = %s", cmdId)
		}
	}

	// Randomly chooses a command if no command has a majority of votes
	for _, cmd := range cons.cmdMap {
		return cmd
	}

	logger.Fatalf("Cannot choose a command. Debug! Accepted cmdMap = %v voteMap = %v",
		cons.cmdMap, cons.cmdVoteMap)
	return nil
}

func (cons *FpCons) IsSetFpCmdConsRet(cmdId string) bool {
	if _, ok := cons.acceptCmdConsRetMap[cmdId]; !ok {
		logger.Fatalf("Cannot find cmd consensus result flag for cmdId = %s", cmdId)
	}
	return cons.acceptCmdConsRetMap[cmdId]
}

func (cons *FpCons) SetFpCmdConsRet(cmdId string) {
	if set, ok := cons.acceptCmdConsRetMap[cmdId]; !ok {
		logger.Fatalf("Cannot find cmd consensus result flag for cmdId = %s", cmdId)
	} else if set {
		logger.Fatalf("Cannot set cmd consenus result flag twice cmdId = %s", cmdId)
	}

	cons.acceptCmdConsRetMap[cmdId] = true
}

func (cons *FpCons) GetRejectCmdIdList() []string {
	ret := make([]string, 0, len(cons.acceptCmdConsRetMap))
	for cmdId, set := range cons.acceptCmdConsRetMap {
		if !set {
			ret = append(ret, cmdId)
		}
	}
	return ret
}

func (cons *FpCons) GetCmdIdMap() map[string]bool {
	ret := make(map[string]bool)
	for cmdId, _ := range cons.acceptCmdConsRetMap {
		ret[cmdId] = true
	}
	return ret
}
