package storing

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
