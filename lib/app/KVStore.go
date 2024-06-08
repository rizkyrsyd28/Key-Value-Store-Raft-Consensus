package app

import "errors"

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

func (kv *KVStore) Get(key string) (string, error) {
	value, ok := kv.store[key]
	if !ok {
		return "", errors.New("key not found")
	}
	return value, nil
}

func (kv *KVStore) Set(key, value string) (string, error) {
	if key == "" {
		return "ERROR", errors.New("key cannot be empty")
	}
	kv.store[key] = value
	return "OK", nil
}

func (kv *KVStore) Strlen(key string) (int, error) {
	value, err := kv.Get(key)
	if err != nil {
		return 0, err
	}
	return len(value), nil
}

func (kv *KVStore) Delete(key string) (string, error) {
	value, ok := kv.store[key]
	if !ok {
		return "", errors.New("key not found")
	}
	delete(kv.store, key)
	return value, nil
}

func (kv *KVStore) Append(key, value string) (string, error) {
	existing, err := kv.Get(key)
	if err != nil {
		return "", err
	}
	new_value := existing + value
	kv.store[key] = new_value
	return "OK", nil
}
