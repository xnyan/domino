package kv

import (
	"fmt"
)

// Not thread-safe
type KvStore struct {
	kvMap map[string]string
}

func NewKvStore() *KvStore {
	return &KvStore{
		kvMap: make(map[string]string),
	}
}

func (kv *KvStore) Write(key, val string) {
	kv.kvMap[key] = val
}

func (kv *KvStore) Read(key string) (string, bool) {
	val, exists := kv.kvMap[key]

	return val, exists
}

func (kv *KvStore) InitData(keyList []string, initVal string) {
	for _, k := range keyList {
		kv.kvMap[k] = initVal
	}
}

// Testing function
func (kv *KvStore) Dump() {
	for key, val := range kv.kvMap {
		fmt.Println(key + " = " + val)
	}
}
