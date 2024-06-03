package util

type ClusterNode struct {
	Address Address
	AckLn   int64
	SentLn  int64
}
