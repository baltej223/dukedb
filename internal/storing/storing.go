// Package storing provides the core key-value storage engine.
package storing

import "sync"

type KV struct {
	mu sync.RWMutex
	m  map[string][]byte
}

var KVGlobalStruct KV

func InitialiseKV() *KV {
	m := make(map[string][]byte)
	KVGlobalStruct.m = m
	return &KVGlobalStruct
}

func Put(key string, value []byte) {
	KVGlobalStruct.mu.Lock()
	defer KVGlobalStruct.mu.Unlock()

	KVGlobalStruct.m[key] = value
}

func Get(key string) ([]byte, bool) {
	KVGlobalStruct.mu.RLock()
	defer KVGlobalStruct.mu.RUnlock()

	value, ok := KVGlobalStruct.m[key]
	return value, ok
}

func Delete(key string) {
	KVGlobalStruct.mu.Lock()
	defer KVGlobalStruct.mu.Unlock()

	delete(KVGlobalStruct.m, key)
}

func Exists(key string) bool {
	KVGlobalStruct.mu.RLock()
	defer KVGlobalStruct.mu.RUnlock()

	_, ok := KVGlobalStruct.m[key]
	return ok
}

func Count() int {
	KVGlobalStruct.mu.RLock()
	defer KVGlobalStruct.mu.RUnlock()

	return len(KVGlobalStruct.m)
}
