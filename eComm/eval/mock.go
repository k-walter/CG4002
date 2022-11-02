package eval

import (
	"cg4002/eComm/common"
	pb "cg4002/protos"
)

type Mock struct {
	*pb.State
}

func MakeMock(args *common.Arg) *Mock {
	return &Mock{}
}

func (c *Mock) Run() {
}

func (c *Mock) Close() {
}

func (c *Mock) BlockingSend(s *pb.State) {
	c.State = s
}

func (c *Mock) BlockingRecv() *pb.State {
	return c.State
}
