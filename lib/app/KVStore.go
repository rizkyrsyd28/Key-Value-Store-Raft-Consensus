package app

import (
	"errors"
	"strconv"
	"strings"
	"sync"
)

// Command Type
const (
	PING    = "PING"
	SET     = "SET"
	GET     = "GET"
	DELETE  = "DELETE"
	APPEND  = "APPEND"
	STRLEN  = "STRLEN"
)

type KVStore struct {
	store map[string]string
	lock  sync.Mutex
}

func NewKVStore() *KVStore {
	return &KVStore{
		store: make(map[string]string),
	}
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
	
	case "REQUEST_LOG":
		// TODO: implement request log
		return "", errors.New("REQUEST_LOG is not implemented")

	default:
		return "", errors.New("unknown command")
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
	kv.lock.Lock()
	defer kv.lock.Unlock()
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
	kv.lock.Lock()
	defer kv.lock.Unlock()
	value, ok := kv.store[key]
	if !ok {
		return "", errors.New("key did not find any value")
	}
	delete(kv.store, key)
	return value, nil
}

func (kv *KVStore) Append(key, value string) (string, error) {
	existing, err := kv.Get(key)
	if err != nil {
		return "", err
	}
	kv.lock.Lock()
	defer kv.lock.Unlock()
	new_value := existing + value
	kv.store[key] = new_value
	return "OK", nil
}

func IsCommandIdempotent(commandType string) bool {
	if( commandType != PING && commandType != GET && commandType != STRLEN) {
		return false
	}
	return true
}

func GetCommandType(command string) string {
	command = strings.ToUpper(strings.Fields(command)[0])
	switch command {
	case PING:
		return PING
	case SET:
		return SET
	case GET:
		return GET
	case DELETE:
		return DELETE
	case APPEND:
		return APPEND
	case STRLEN:
		return STRLEN
	default:
		return ""
	}
}