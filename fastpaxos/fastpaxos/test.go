package fastpaxos

func (f *FastPaxos) Test() {
	logger.Infof("Testing FastPaxos: log size = %s, nextEmptyIdx = %s, nextExecIdx = %s, nextPersistIdx = %s",
		f.lm.GetLogSize(), f.lm.GetNextEmptyIdx(), f.lm.GetNextExecIdx(), f.lm.GetNextPersistIdx())

	f.fpManager.Test()
}

func (fp *FastPathManager) Test() {
	logger.Infof("Testing FastPath Manager (TFM):")
	logger.Infof("(TFM): fastRetTable size = %d", len(fp.opFastRetTable))
	logger.Infof("(TFM): slowRetTable size = %d", len(fp.opSlowRetTable))

	logger.Infof("(TFM): opTable size = %d", len(fp.opTable))
	for opId, opInfo := range fp.opTable {
		logger.Infof("opId %s => %s", opId, opInfo)
	}

	logger.Infof("(TFM): idxTable size = %d", len(fp.idxTable))
	for idx, idxInfo := range fp.idxTable {
		logger.Infof("idx %s => %s", idx, idxInfo)
	}

	logger.Infof("(TFM): opQueue size = %d", fp.opQueue.Size())
	logger.Infof("(TFM): opQueue: %s", fp.opQueue)
	logger.Infof("(TFM): idxQueue size = %d", fp.idxQueue.Size())
	logger.Infof("(TFM): idxQueue: %s", fp.idxQueue)
}
