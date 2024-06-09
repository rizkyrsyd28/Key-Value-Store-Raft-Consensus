package app

import (
	"errors"
	"strconv"
	"strings"

)

type KVStore struct {
	store map[string]string
}

func (kv *KVStore) Execute(commandRaw string) (string, error) {
    // Split the command string into parts
    parts := strings.Fields(commandRaw)
    if len(parts) == 0 {
        return "", errors.New("empty command")
    }

    command := strings.ToUpper(parts[0])
    switch command {
	case "PING":
		return kv.Ping(), nil
    case "SET":
        if len(parts) != 3 {
            return "", errors.New("invalid SET command format, expected: SET key value")
        }
        key, value := parts[1], parts[2]
        value, err := kv.Set(key, value)
		if err != nil {
			return value, err
		}
        return value, nil

    case "GET":
        if len(parts) != 2 {
            return "", errors.New("invalid GET command format, expected: GET key")
        }
        key := parts[1]
        value, err := kv.Get(key)
        if err != nil {
            return value, err
        }
        return value, nil

    case "DELETE":
        if len(parts) != 2 {
            return "", errors.New("invalid DELETE command format, expected: DELETE key")
        }
        key := parts[1]
        value, err := kv.Delete(key)
		if err != nil {
			return value, err
		}
        return value, nil
	
	case "APPEND":
		if len(parts) != 3 {
			return "", errors.New("invalid APPEND command format, expected: APPEND key value")
		}
		key, value := parts[1], parts[2]
		value, err := kv.Append(key, value)
		if err != nil {
			return value, err
		}
		return value, nil
		
	case "STRLEN":
		if len(parts) != 2 {
			return "", errors.New("invalid STRLEN command format, expected: STRLEN key")
		}
		key := parts[1]
		value, err := kv.Strlen(key)
		if err != nil {
			return "", err
		}
		return strconv.Itoa(value), nil

    default:
        return "", errors.New("unknown command")
    }
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
