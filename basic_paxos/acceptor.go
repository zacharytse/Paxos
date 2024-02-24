package basic_paxos

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

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

// Prepare 只有提议者发送来的提案编号大于接受者见过的最大提案编号才会接受
// 如果接受者之前有接收别的提案，则响应消息还应包含已接收的提案信息/**
func (a *Acceptor) Prepare(args *MsgArgs, reply *MsgReply) error {
	if args.Number > a.minProposal {

		a.minProposal = args.Number
		reply.Number = a.acceptedNumber
		reply.Value = a.acceptedValue
		reply.Ok = true
	} else {
		reply.Ok = false
	}
	return nil
}

// Accept 接受者没有收到比最大提案编号更大的提案则接收该提案，
// 同时将该提案转发给全部学习者/**
func (a *Acceptor) Accept(args *MsgArgs, reply *MsgReply) error {
	if args.Number >= a.minProposal {
		a.minProposal = args.Number
		a.acceptedNumber = args.Number
		a.acceptedValue = args.Value
		reply.Ok = true
		// 后台转发接收的提案给学习者
		for _, lid := range a.learners {
			go func(learner int) {
				addr := fmt.Sprintf("127.0.0.1:%d", learner)
				args.From = a.id
				args.To = learner
				resp := new(MsgReply)
				ok := call(addr, "Learner.Learn", args, resp)
				if !ok {
					return
				}
			}(lid)
		}
	} else {
		reply.Ok = false
	}
	return nil
}

// 初始化接受者并启动服务
func newAcceptor(id int, learners []int) *Acceptor {
	acceptor := &Acceptor{
		id:       id,
		learners: learners,
	}
	acceptor.server()
	return acceptor
}

func (a *Acceptor) server() {
	rpcs := rpc.NewServer()
	rpcs.Register(a) // 绑定Accept和Prepare方法
	addr := fmt.Sprintf(":%d", a.id)
	l, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	a.lis = l
	go func() {
		for {
			conn, err := a.lis.Accept()
			if err != nil {
				continue
			}
			go rpcs.ServeConn(conn)
		}
	}()
}

func (a *Acceptor) close() {
	a.lis.Close()
}
