package main

import (
	"time"
)

type State uint

const (
	FOLLOWER State = iota
	CANDIDATE
	LEADER
)

type Server struct {
	state            State
	term             uint
	electionTimeout  Timeout
	heartbeatTimeout Timeout
}

func main() {
	server := &Server{
		state:            FOLLOWER,
		term:             0,
		electionTimeout:  *createRandomTimeout(150, 300, time.Millisecond),
		heartbeatTimeout: *createRandomTimeout(150, 300, time.Millisecond),
	}

	for {
		switch server.state {
		case FOLLOWER:
			listenQueues(server)
			break
		case CANDIDATE:
			break
		case LEADER:
			LBCommunication(server)
			break
		}
	}
}
