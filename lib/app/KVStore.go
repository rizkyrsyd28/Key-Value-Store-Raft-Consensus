package app

type KVStore struct {
	store map[string]string
}

func NewKVStore() *KVStore {
	return &KVStore{
		store: make(map[string]string),
	}
}

func (kv *KVStore) Ping() string {
	return "PONG"
}

func (kv *KVStore) Get(key string) string {
	value, ok := kv.store[key]
	if !ok {
		return ""
	}
	return value
}

func (kv *KVStore) Set(key, value string) string {
	kv.store[key] = value
	return "OK"
}

func (kv *KVStore) Strlen(key string) int {
	value := kv.Get(key)
	return len(value)
}

func (kv *KVStore) Delete(key string) string {
	value, ok := kv.store[key]
	if !ok {
		return ""
	}
	delete(kv.store, key)
	return value
}

func (kv *KVStore) Append(key, value string) string {
	existing := kv.Get(key)
	new_value := existing + value
	kv.store[key] = new_value
	return "OK"
}
