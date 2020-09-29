package clientlib

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"domino/common"
	"domino/dynamic/latency"
)

func TestDynamicLinearLat(t *testing.T) {
	configFile := "./testdata/dynamic_test.config"
	lib := newDynamicClientLib("Test", configFile)
	lib.startProcessing()

	p := common.NewProperties()
	p.Load(configFile)
	rList := lib.replicaAddrList

	fmt.Println("====Linear Latency Test====")

	wt, _ := time.ParseDuration("5ms")
	for k := 0; k < 100; k++ {
		for i, r := range rList {
			t, _ := time.ParseDuration("10ms")
			t = t*time.Duration(i+1) + time.Duration(k*1000000)
			ret := &latency.ProbeRet{r, t}
			lib.probeC <- ret
			fmt.Println("Probe", r, t.Milliseconds())
		}
		fmt.Println(lib.selectFp())
		time.Sleep(wt)
	}
}

func TestDynamicRandomLat(t *testing.T) {
	configFile := "./testdata/dynamic_test.config"
	lib := newDynamicClientLib("Test", configFile)
	lib.startProcessing()

	p := common.NewProperties()
	p.Load(configFile)
	rList := lib.replicaAddrList

	fmt.Println("====Random Latency Test====")

	wt, _ := time.ParseDuration("5ms")
	for k := 0; k < 10; k++ {
		for _, r := range rList {
			t, _ := time.ParseDuration("10ms")
			t = t * time.Duration(rand.Int31n(5))
			ret := &latency.ProbeRet{r, t}
			lib.probeC <- ret
			fmt.Println("Probe", r, t.Milliseconds())
		}
		fmt.Println(lib.selectFp())
		time.Sleep(wt)
	}
}
