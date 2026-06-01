package storing

import "encoding/json"

var KVInternalStruct KV

func InitialiseKVI() *KV {
	m := make(map[string][]byte)
	KVInternalStruct.m = m
	return &KVInternalStruct
}

func PutI(key string, value []byte) {
	KVInternalStruct.mu.Lock()
	defer KVInternalStruct.mu.Unlock()

	KVInternalStruct.m[key] = value
}

func PutIJSON(key string, value any) {
	bytes, err := json.Marshal(value)
	if err != nil {
		panic(err) // or return the error
	}
	PutI(key, bytes)
}

func GetIJSON[T any](key string) (T, bool) {
	var value T

	data, ok := GetI(key)
	if !ok {
		return value, false
	}

	if err := json.Unmarshal(data, &value); err != nil {
		panic(err) // or return the error
	}

	return value, true
}

func GetI(key string) ([]byte, bool) {
	KVInternalStruct.mu.RLock()
	defer KVInternalStruct.mu.RUnlock()

	value, ok := KVInternalStruct.m[key]
	return value, ok
}

func DeleteI(key string) {
	KVInternalStruct.mu.Lock()
	defer KVInternalStruct.mu.Unlock()

	delete(KVInternalStruct.m, key)
}

func ExistsI(key string) bool {
	KVInternalStruct.mu.RLock()
	defer KVInternalStruct.mu.RUnlock()

	_, ok := KVInternalStruct.m[key]
	return ok
}

func CountI() int {
	KVInternalStruct.mu.RLock()
	defer KVInternalStruct.mu.RUnlock()

	return len(KVInternalStruct.m)
}
