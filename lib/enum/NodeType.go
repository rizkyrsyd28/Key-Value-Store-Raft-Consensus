package enum

type NodeType string

const (
	LEADER    NodeType = "LEADER"
	CANDIDATE NodeType = "CANDIDATE"
	FOLLOWER  NodeType = "FOLLOWER"
)
