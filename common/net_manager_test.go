package common

import (
	"testing"
	"time"
)

func TestStaticNetworkManager(t *testing.T) {
	return
	// TODO change the loadin file to a fixed local file
	configFile := "../config/delay-conf.json.template"
	var m NetworkManager = NewStaticNetworkManager(configFile, "oneway-delay")

	if m.GetOneWayNetDelay("D1", "D1") != time.Duration(0.2*1000*1000) {
		t.Errorf("D1 -> D1 != 0.2ms")
	}
	if m.GetOneWayNetDelay("D1", "D2") != time.Duration(50*1000*1000) {
		t.Errorf("D1 -> D2 != 50ms")
	}
	if m.GetOneWayNetDelay("D6", "D3") != time.Duration(90*1000*1000) {
		t.Errorf("D6 -> D3 != 90ms")
	}

	dcList := []string{"D1", "D2", "D3"}
	if m.MaxOneWayNetDelay("D1", dcList) != time.Duration(60*1000*1000) {
		t.Errorf("max(D1-->{D1, D2, D3}) != 60ms")
	}
	if m.MaxOneWayNetDelay("D4", dcList) != time.Duration(70*1000*1000) {
		t.Errorf("max(D4-->{D1, D2, D3}) != 70ms")
	}

	l := m.GetClosestQuorum("D1", dcList, 1)
	if len(l) != 1 {
		t.Errorf("D1 closest DCs: %v", l)
	}
	if l[0] != "D1" {
		t.Errorf("D1 closest DC is not D1 but %v %d", l, len(l))
	}
	l = m.GetClosestQuorum("D3", dcList, 1)
	if len(l) != 1 {
		t.Errorf("D3 closest DCs: %v", l)
	}
	if l[0] != "D3" {
		t.Errorf("D3 closest DC is not D3 but %v %d", l, len(l))
	}
	l = m.GetClosestQuorum("D6", dcList, 1)
	if len(l) != 1 {
		t.Errorf("D6 closest DCs: %v", l)
	}
	if l[0] != "D1" {
		t.Errorf("D6 closest DC is not D1 but %v %d", l, len(l))
	}

	if m.MaxDifferenceNetDelay("D1", dcList) != time.Duration((60-0.2)*1000*1000) {
		t.Errorf("maxDiff(D1-->{D1, D2, D3}) != 59.8ms")
	}
	if m.MaxDifferenceNetDelay("D2", dcList) != time.Duration((60-0.2)*1000*1000) {
		t.Errorf("maxDiff(D2-->{D1, D2, D3}) != 59.8ms")
	}
	if m.MaxDifferenceNetDelay("D5", dcList) != time.Duration(30*1000*1000) {
		t.Errorf("maxDiff(D2-->{D1, D2, D3}) != 30ms")
	}
}
