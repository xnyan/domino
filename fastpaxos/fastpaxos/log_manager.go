package fastpaxos

import (
	"sync"
	"time"

	"domino/fastpaxos/rpc"
)

type LogManager interface {
	// Fast path
	FastPathAccept(entry *Entry) (*LogIdx, error) // accepts the entry to an empty log entry

	// Slow path
	SlowPathAccept(idx *LogIdx, entry *Entry) error

	Commit(idx *LogIdx, entry *Entry) error

	// Output for execution
	Exec(execCh chan *rpc.Operation) // log execution output channel

	// Log persistence
	Persist()

	// Start
	Run(execCh chan *rpc.Operation)

	// Helper functions for testing
	GetLogSize() string
	GetNextEmptyIdx() *LogIdx
	GetNextExecIdx() *LogIdx
	GetNextPersistIdx() *LogIdx
}

type DefaultLogManager struct {
	log Log // Log

	logNextEmptyIdx *LogIdx

	// Execution
	logNextExecIdx *LogIdx // next entry index to execute in the log
	logExecCond    *sync.Cond

	// Log persistence
	isAsyncPersist    bool
	pLogger           *Logger
	logNextPersistIdx *LogIdx
	logPersistCond    *sync.Cond
}

func NewDefaultLogManager(
	l Log,
	isAsyncP bool,
	logger *Logger,
) *DefaultLogManager {
	m := &DefaultLogManager{
		log:             l,
		logNextEmptyIdx: &LogIdx{0, 0},

		logNextExecIdx: &LogIdx{0, 0},

		isAsyncPersist:    isAsyncP,
		pLogger:           logger,
		logNextPersistIdx: &LogIdx{0, 0},
	}

	var execLock sync.Mutex
	m.logExecCond = sync.NewCond(&execLock)

	var persistLock sync.Mutex
	m.logPersistCond = sync.NewCond(&persistLock)

	return m
}

func (m *DefaultLogManager) FastPathAccept(entry *Entry) (*LogIdx, error) {
	for true {
		e, err := m.log.Get(m.logNextEmptyIdx)
		if err != nil {
			logger.Fatalf("FastPathAccept() gets log entry error: %v", err)
		}
		if e != nil {
			// It is possible that the entry has been assigned with an operation via the slow path

			logger.Debugf("Fast-path cannot accept opId = (%s) at idx (%s) [opId (%s), status (%d)]",
				entry.op.Id, m.logNextEmptyIdx, e.op.Id, e.status)

			if err = m.log.IncIdx(m.logNextEmptyIdx); err != nil {
				logger.Fatalf("FastPathAccept() (lookup) fails to increase idx, error %v", err)
			}
			continue
		}
		break
	}

	err := m.log.Put(m.logNextEmptyIdx, entry)
	if err != nil {
		logger.Fatalf("FastPasthAccept() writes opId = (%s) at log idx = (%s) fails. error: %v",
			entry.op.Id, m.logNextEmptyIdx, err)
	}

	idx := &LogIdx{m.logNextEmptyIdx.Seg, m.logNextEmptyIdx.Offset}
	if err = m.log.IncIdx(m.logNextEmptyIdx); err != nil {
		logger.Fatalf("FastPathAccept() (next) fails to increase idx, error %v", err)
	}

	return idx, nil
}

func (m *DefaultLogManager) SlowPathAccept(idx *LogIdx, entry *Entry) error {
	e, err := m.log.Get(idx)
	if err != nil {
		logger.Fatalf("SlowPathAccept() gets log entry error: %v", err)
	}

	if e == nil {
		// empty log entry
		err = m.log.Put(idx, entry)
		if err != nil {
			logger.Fatalf("SlowPasthAccept() fails to put opId = (%s) at an (empty) log idx = (%s)"+
				", error: %v", entry.op.Id, idx, err)
		}
		entry.SetSlowAccepted()

		logger.Debugf("Slow-path accepted opId = (%s) at idx = (%s)", entry.op.Id, idx)

		return nil
	}

	if e.IsCommitted() {
		if e.op.Id != entry.op.Id {
			logger.Fatalf("Slow-path cannot accept operation id = (%s) at "+
				"a committed log idx = (%s) [opId (%s)]", entry.op.Id, idx, e.op.Id)
		}
	} else if e.IsSlowAccepted() {
		logger.Fatalf("Slow-path cannot accept operation id = (%s) at "+
			"a slow-path accepted log idx = (%s) [opId %s]", entry.op.Id, idx, e.op.Id)
	} else {
		if e.op.Id != entry.op.Id {
			err = m.log.Put(idx, entry)
			if err != nil {
				logger.Fatalf("SlowPasthAccept() fails to put opId = (%s) at log idx = (%s)"+
					", error: %v", entry.op.Id, idx, err)
			}
		}
		entry.SetSlowAccepted()

		logger.Debugf("Slow-path accepted opId = (%s) at idx = (%s)", entry.op.Id, idx)
	}

	return nil
}

func (m *DefaultLogManager) Commit(idx *LogIdx, entry *Entry) error {
	if m.isAsyncPersist {
		defer m.logPersistCond.Signal()
	}
	// Uses a cv to notify execution thread to execute committed log entries. Do
	// not use a channel here to avoid blocking on the sender side.
	defer m.logExecCond.Signal()

	e, err := m.log.Get(idx)
	if err != nil {
		logger.Fatalf("Commit() gets log entry error: %v", err)
	}

	// It is possible that a replica receives the leader's commit request before the proposal
	if e == nil {
		err = m.log.Put(idx, entry)
		if err != nil {
			logger.Fatalf("Commit() fails to put opId = (%s) at an (empty) log idx = (%s)"+
				", error: %v", entry.op.Id, idx, err)
		}
		entry.SetCommitted()
		logger.Debugf("Commit() committed opId = (%s) at empty idx = (%s)", entry.op.Id, idx)

		return nil
	}

	if e.IsCommitted() {
		logger.Fatalf("Commit() cannot commit opId = (%s) at "+
			"a committed log idx = (%s) [opId = (%s)]", entry.op.Id, idx, e.op.Id)
	}
	if e.op.Id != entry.op.Id {
		err = m.log.Put(idx, entry)
		if err != nil {
			logger.Fatalf("Commit() fails to put opId = (%s) at log idx = (%s)"+
				", error: %v", entry.op.Id, idx, err)
		}
		entry.SetCommitted()
	} else {
		e.SetCommitted()
	}
	e.SetStartDuration(time.Now().UnixNano())
	logger.Debugf("Commit() committed opId = (%s) at non-empty idx = (%s)", entry.op.Id, idx)

	return nil
}

func (m *DefaultLogManager) Exec(execCh chan *rpc.Operation) {
	logger.Debugf("Starts execution channel")
	// Keeps putting committed operations into a chennel for applications to execute
	for {
		m.logExecCond.L.Lock()
		entry, err := m.log.Get(m.logNextExecIdx)
		for err == nil && (entry == nil || !entry.IsCommitted()) {
			m.logExecCond.Wait()
			entry, err = m.log.Get(m.logNextExecIdx)
		}
		m.logExecCond.L.Unlock()

		if err != nil {
			logger.Fatalf("Exec() error: %v", err)
		}

		logger.Debugf("duration = %v ns",time.Now().UnixNano()-entry.GetStartDuration()) 
		execCh <- entry.GetOp() // May block if the channel is full

		if err = m.log.IncIdx(m.logNextExecIdx); err != nil {
			logger.Fatalf("Exec() fails to increase idx, error: %v ", err)
		}
	}
}

func (m *DefaultLogManager) Persist() {
	logger.Debugf("Starts persistence thread")
	for {
		m.logPersistCond.L.Lock()
		entry, err := m.log.Get(m.logNextPersistIdx)
		for err == nil && (entry == nil || !entry.IsCommitted()) {
			m.logPersistCond.Wait()
			entry, err = m.log.Get(m.logNextPersistIdx)
		}
		m.logPersistCond.L.Unlock()

		if err != nil {
			logger.Fatalf("Persist() error: %v", err)
		}

		logger.Debugf("Persists idx = %s", m.logNextPersistIdx)

		m.pLogger.Write(entry.String() + "\n")

		if err = m.log.IncIdx(m.logNextPersistIdx); err != nil {
			logger.Fatalf("Persist() fails to increase idx, error: %v ", err)
		}
	}
	m.pLogger.Close()
}

func (m *DefaultLogManager) Run(execCh chan *rpc.Operation) {
	go m.Exec(execCh)
	if m.isAsyncPersist {
		go m.Persist()
	}
}

func (m *DefaultLogManager) GetLogSize() string {
	return m.log.Size()
}

func (m *DefaultLogManager) GetNextEmptyIdx() *LogIdx {
	return m.logNextEmptyIdx
}

func (m *DefaultLogManager) GetNextExecIdx() *LogIdx {
	return m.logNextExecIdx
}

func (m *DefaultLogManager) GetNextPersistIdx() *LogIdx {
	return m.logNextPersistIdx
}
