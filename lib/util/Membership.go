package util

type MembershipApply struct {
	Address Address
	Insert  bool
}

func NewMembershipApply(address Address, insert bool) *MembershipApply {
	return &MembershipApply{
		Address: address,
		Insert:  insert,
	}
}
