package BasicPaxos

import "net/rpc"

type MsgArgs struct {
	// 提案编号
	Number int
	// 提案值
	Value interface{}
	// 发送者id
	From int
	// 接受者id
	To int
}

type MsgReply struct {
	Ok     bool
	Number int
	Value  interface{}
}

// 发送rpc消息的函数
func call(srv string, name string, args interface{}, reply interface{}) bool {
	c, err := rpc.Dial("tcp", srv)
	if err != nil {
		return false
	}
	defer c.Close()

	err = c.Call(name, args, reply)
	if err == nil {
		return true
	}
	return false
}
