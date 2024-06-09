package util

import "github.com/Sister20/if3230-tubes-dark-syster/lib/pb"

type ClusterNode struct {
	Address *Address
	AckLn   int64
	SentLn  int64
}

type ClusterNodeList struct {
	Map map[string]ClusterNode
}

func (c *ClusterNodeList) AddAddress(address *Address) {
	c.Map[address.ToString()] = ClusterNode{
		Address: address,
		AckLn:   0,
		SentLn:  0,
	}
}

func (c *ClusterNodeList) SetAddressPb(address []*pb.Address) {
	c.Map = map[string]ClusterNode{}
	for _, member := range address {
		c.AddAddress(&Address{member})
	}
}

func (c *ClusterNodeList) SetAddress(address []*Address) {
	c.Map = map[string]ClusterNode{}
	for _, member := range address {
		c.AddAddress(member)
	}
}

func (c *ClusterNodeList) PatchAddress(address *Address, clusterNode ClusterNode) ClusterNode {
	c.Map[address.ToString()] = clusterNode
	return c.Map[address.ToString()]
}

func (c *ClusterNodeList) Get(id string) ClusterNode {
	return c.Map[id]
}

func (c *ClusterNodeList) GetAllAddress() []Address {
	var addrs []Address
	for _, node := range c.Map {
		addrs = append(addrs, *node.Address)
	}
	return addrs
}

func (c *ClusterNodeList) GetAllPbAddress() []*pb.Address {
	var addrs []*pb.Address
	for _, node := range c.Map {
		addrs = append(addrs, node.Address.Address)
	}
	return addrs
}
