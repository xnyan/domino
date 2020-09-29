package clientlib

import (
	"github.com/op/go-logging"

	"domino/common"
	//"domino/dynamic/config"
	"domino/dynamic/dynamic"
)

var logger = logging.MustGetLogger("DynamicClientLib")

type ClientLib interface {
	/* Chooses a Paxos or Fast Paxos instance to propose the command
	 * Returns (isFastUsed, isCommitted, isFast, execution result)
	 */
	Propose(cmd *dynamic.Command) (bool, bool, bool, string)

	/* Uses an Fast Paxos instance to propose the command
	 * Returns (isCommitted, isFast, execution result)
	 */
	FpPropose(cmd *dynamic.Command) (bool, bool, string)

	/* Uses a Paxos instance to propose the command
	 * Returns (isCommitted, execution result)
	 */
	PaxosPropose(cmd *dynamic.Command) (bool, string)

	// Shuts down the client lib
	Close()

	// Static network configuration
	Test()
}

func NewClientLib(id, dcId, configFile, replicaFile, targetReplicaDcId string) ClientLib {
	p := common.NewProperties()
	p.Load(configFile)

	// Dynamic Predicting Latency
	lib := newDynamicClientLib(id, dcId, configFile, replicaFile, targetReplicaDcId)
	lib.start(false, lib.probeInv)
	return lib

	/*
		mode := p.GetWithDefault(config.FLAG_DYNAMIC_LATENCY_PREDICTION_MODE, "dynamic")
		if mode == "dynamic" {
			lib := newDynamicClientLib(id, configFile)
			lib.start(false, lib.probeInv)
			return lib
		} else if mode == "static" {
			return newStaticClientLib(id, dcId, configFile)
		} else {
			logger.Fatalf("Invalid prediction mode %s, expected dynamic or static", mode)
		}
		return nil
	*/
}
