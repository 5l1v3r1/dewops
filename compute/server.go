package main

import (
	"fmt"
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
	fmt.Println(server.state)
	fmt.Println(server.term)
	fmt.Println(server.electionTimeout.period)
	fmt.Println(server.heartbeatTimeout.period)
}
