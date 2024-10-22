package persistent_storage

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/Sister20/if3230-tubes-dark-syster/lib/logger"
	"github.com/Sister20/if3230-tubes-dark-syster/lib/util"
)

type PersistValues struct {
	ElectionTerm    uint64
	VotedFor        *util.Address
	Log             logger.RaftNodeLog
	CommittedLength uint64
}

type PersistentStorage struct {
	ID   string
	Path string
	Lock sync.Mutex
}

func NewPersistentStorage(addr *util.Address) *PersistentStorage {
	id := fmt.Sprintf("%s_%d", addr.IP, addr.Port)
	path := fmt.Sprintf("persistent/%s.json", id)
	// make folder if not exist
	if _, err := os.Stat("persistent"); os.IsNotExist(err) {
		os.Mkdir("persistent", 0755)
	}

	return &PersistentStorage{
		ID:   id,
		Path: path,
	}
}

func (ss *PersistentStorage) Load() *PersistValues {
	ss.Lock.Lock()
	defer ss.Lock.Unlock()

	data, err := os.ReadFile(ss.Path)
	if err != nil {
		return nil
	}

	var result PersistValues
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil
	}

	return &result
}

func (ss *PersistentStorage) StoreAll(data *PersistValues) {
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

func (ss *PersistentStorage) InitialStoreAll() error {
	initialData := `{"ElectionTerm":1,"VotedFor":null,"Log":{"entries":[{"term":1,"command":"PING"}]},"CommittedLength":0}`

	err := os.WriteFile(ss.Path, []byte(initialData), 0644)
	if err != nil {
		return err
	}

	return nil
}
