package enum

type RaftConst float32

const (
	HEARTBEAT_INTERVAL   RaftConst = 1
	ELECTION_TIMEOUT_MIN RaftConst = 2
	ELECTION_TIMEOUT_MAX RaftConst = 3
	RPC_TIMEOUT          RaftConst = 0.5
)
