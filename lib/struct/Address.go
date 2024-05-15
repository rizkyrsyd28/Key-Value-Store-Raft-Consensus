package _struct

import "fmt"

type Address struct {
	IP   string
	Port int
}

func New(ip string, port int) *Address {
	return &Address{IP: ip, Port: port}
}

func (address Address) ToString() string {
	return fmt.Sprintf("%s:%d", address.IP, address.Port)
}

func (address Address) Iterator() <-chan interface{} {
	iter := make(chan interface{}, 2)
	iter <- address.IP
	iter <- address.Port
	close(iter)
	return iter
}

func (address Address) IsEqual(other Address) bool {
	return address.IP == other.IP && address.Port == other.Port
}

func (address Address) IsNotEqual(other Address) bool {
	return address.IP != other.IP || address.Port != other.Port
}
