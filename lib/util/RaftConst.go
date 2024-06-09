package util

type RaftConst float32

const (
	HEARTBEAT_INTERVAL   int = 4
	ELECTION_TIMEOUT_MIN int = 15
	ELECTION_TIMEOUT_MAX int = 30
	RPC_TIMEOUT          RaftConst = 0.5
	IDLE                 RaftConst = 0
)
