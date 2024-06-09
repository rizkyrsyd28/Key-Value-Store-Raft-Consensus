package stable_storage

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/Sister20/if3230-tubes-dark-syster/lib/logger"
)

type Address struct {
	IP   string
	Port int
}

type StableVars struct {
	ElectionTerm uint64
	VotedFor     *Address
	Log          logger.RaftNodeLog
	CommitLength uint64
}

type StableStorage struct {
	ID   string
	Path string
	Lock sync.Mutex
}

func NewStableStorage(addr Address) *StableStorage {
	id := fmt.Sprintf("%s_%d", addr.IP, addr.Port)
	path := fmt.Sprintf("/persistence/%s.json", id)
	fmt.Println("id", id)
	fmt.Println("path", path)

	return &StableStorage{
		ID:   id,
		Path: path,
	}
}

func (ss *StableStorage) Load() *StableVars {
	ss.Lock.Lock()
	defer ss.Lock.Unlock()

	data, err := os.ReadFile(ss.Path)
	if err != nil {
		return nil
	}

	var result StableVars
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil
	}

	return &result
}

func (ss *StableStorage) StoreAll(data *StableVars) {
	ss.Lock.Lock()
	defer ss.Lock.Unlock()

	strData, err := json.Marshal(data)
	if err != nil {
		return
	}

	err = os.WriteFile(ss.Path, strData, 0644)
	if err != nil {
		return
	}
}

func (ss *StableStorage) TryLoad() *StableVars {
	ss.Lock.Lock()
	defer ss.Lock.Unlock()

	return ss.Load()
}
