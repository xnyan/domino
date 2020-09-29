package main

type Client interface {
	// Returns isFastUsed, isAccept, isFast, ExecRet,
	ExecTxn(rKeyList []string, wTable map[string]string) (bool, bool, bool, string)
	//thread-safe
	SyncExecTxn(rKeyList []string, wTable map[string]string) (bool, bool, bool, string)
	Close()
}
