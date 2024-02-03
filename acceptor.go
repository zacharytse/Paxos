package BasicPaxos

import "net"

type Acceptor struct {
	lis net.Listener
	// 服务器id
	id int
	// 接受者承诺的提案编号，为0表示没有收到过Prepare消息
	minProposal int
	// 接受者已接收的提案编号，为0则表示没有接受任何提案
	acceptedNumber int
	// 接受者已接收的提案值，如果没有接受任何提案，则为nil
	acceptedValue interface{}

	// 学习者id列表
	learners []int
}
