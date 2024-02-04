package BasicPaxos

import (
	"testing"
)

func start(accceptorIds []int, learnerIds []int) ([]*Acceptor, []*Learner) {
	acceptors := make([]*Acceptor, 0)
	for _, aid := range accceptorIds {
		a := newAcceptor(aid, learnerIds)
		acceptors = append(acceptors, a)
	}
	learners := make([]*Learner, 0)
	for _, lid := range learnerIds {
		l := newLearner(lid, accceptorIds)
		learners = append(learners, l)
	}

	return acceptors, learners
}

func cleanup(acceptors []*Acceptor, learners []*Learner) {
	for _, a := range acceptors {
		a.close()
	}

	for _, l := range learners {
		l.close()
	}
}

func TestSingleProposer(t *testing.T) {
	// 1001、1002、1003是接受者id
	acceptorIds := []int{1001, 1002, 1003}
	// 2001是学习者id
	learnerIds := []int{2001}
	acceptors, learns := start(acceptorIds, learnerIds)

	defer cleanup(acceptors, learns)
	// 1是提议者id
	p := &Proposer{
		id:        1,
		acceptors: acceptorIds,
	}

	value := p.propose("hello world")
	if value != "hello world" {
		t.Errorf("value = %s, excepted %s", value, "hello world")
	}

	learnValue := learns[0].chosen()
	if learnValue != value {
		t.Errorf("learnValue = %s,excepted %s", learnValue, "hello world")
	}
}
