package common

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

type NetworkManager interface {
	// Returns the one-way network delay from datacenter d1 to datacenter d2
	GetOneWayNetDelay(dc1, dc2 string) time.Duration

	// Returns the max one-way network delay from dc to any datacenter in dcList
	MaxOneWayNetDelay(dc string, dcList []string) time.Duration

	// Returns the difference between min(dc-->dcList) and max(dc-->dcList)
	MaxDifferenceNetDelay(dc string, dcList []string) time.Duration

	// For dc, returns the closest quorum of datacenters from the dcList
	GetClosestQuorum(dc string, dcList []string, quorum int) []string
}

type StaticNetworkManager struct {
	latMap map[string]map[string]time.Duration // 1-way delay between DCs
}

// The config file should be a json file
func NewStaticNetworkManager(
	configFile string, tag string,
) *StaticNetworkManager {
	m := &StaticNetworkManager{
		latMap: make(map[string]map[string]time.Duration),
	}

	// Parses the json configuration file
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		logger.Fatalf("Cannot read file: %s \nError: %v", configFile, err)
	}
	config := make(map[string]interface{})
	json.Unmarshal(data, &config)

	t := config[tag].(map[string]interface{})

	for dc1, lm := range t {
		m.latMap[dc1] = make(map[string]time.Duration)
		//m.latMap[dc1][dc1], err = time.ParseDuration("0.2ms")// Assumes 0.2ms one-way delay within a dc
		m.latMap[dc1][dc1], err = time.ParseDuration("0ms")
		if err != nil {
			logger.Fatalf("Sets one-way delay within dcId = %s error: %v", dc1, err)
		}

		tm := lm.(map[string]interface{})
		for dc2, lat := range tm {
			m.latMap[dc1][dc2], err = time.ParseDuration(lat.(string))
			if err != nil {
				logger.Fatalf("File = %s, tag = %s, delay from %s to %s error: %v",
					configFile, tag, dc1, dc2, err)
			}
		}
	}

	return m
}

func (m *StaticNetworkManager) GetOneWayNetDelay(dc1, dc2 string) time.Duration {
	if lm, ok := m.latMap[dc1]; ok {
		if l, ok := lm[dc2]; ok {
			return l
		}
	}
	logger.Fatalf("No one-way delay from %s to %s", dc1, dc2)
	return time.Duration(0)
}

func (m *StaticNetworkManager) MaxOneWayNetDelay(dc string, dcList []string) time.Duration {
	lm, ok := m.latMap[dc]
	if !ok {
		logger.Fatalf("No dc = %s", dc)
	}

	var max time.Duration = -1
	for _, dst := range dcList {
		if lat, ok := lm[dst]; ok {
			if lat > max {
				max = lat
			}
		} else {
			logger.Fatalf("No one-way delay from %s to %s", dc, dst)
		}
	}
	return max
}

func (m *StaticNetworkManager) MaxDifferenceNetDelay(dc string, dcList []string) time.Duration {
	max := m.MaxOneWayNetDelay(dc, dcList)
	closestDc := m.GetClosestQuorum(dc, dcList, 1)[0]
	min, ok := m.latMap[dc][closestDc]
	if !ok {
		logger.Fatalf("No one-way delay from %s to %s", dc, closestDc)
	}
	if min > max {
		logger.Fatalf("Min > Max for dc = %s, dcList = %v", dc, dcList)
	}
	return max - min
}

func (m *StaticNetworkManager) GetClosestQuorum(
	dc string, dcList []string, quorum int,
) []string {
	ret := make([]string, 0, len(dcList))
	ret = append(ret, dcList[0:]...)

	k := len(ret) - quorum
	for i := 0; i < k; i++ {
		max := m.GetOneWayNetDelay(dc, ret[i])
		for j := i + 1; j < len(ret); j++ {
			lat := m.GetOneWayNetDelay(dc, ret[j])
			if lat > max {
				ret[i], ret[j] = ret[j], ret[i]
				max = lat
			}
		}
	}
	return ret[k:]
}
