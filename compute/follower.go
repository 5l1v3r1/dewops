package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func listenQueues(server *Server) {
	conn := Connect("amqp://guest:guest@172.18.0.2:5672/")
	// Make sure we close the connection whenever the program is about to exit.
	defer conn.Close()

	ch := GetChannel(conn)
	// Make sure we close the channel whenever the program is about to exit.
	defer ch.Close()

	ip := GetLocalIP()

	// work queue
	work := ip + "_work"
	workMsgs := CreateAndSubscribeQueue(ch, QueueDef{
		Exchange: work,
		Queue:    work,
		Binding:  work,
	}, "topic")
	listenWorkQ(workMsgs)

	// election queue
	elect := ip + "_election"
	electionMsgs := CreateAndSubscribeQueue(ch, QueueDef{
		Exchange: elect,
		Queue:    elect,
		Binding:  elect,
	}, "fanout")
	listenElectionQ(electionMsgs)

	// heartbeat queue
	heart := ip + "_heart"
	hbMsgs := CreateAndSubscribeQueue(ch, QueueDef{
		Exchange: heart,
		Queue:    heart,
		Binding:  heart,
	}, "fanout")
	listenHeartbeatQ(hbMsgs)

	sendVote(server.electionTimeout)

	forever := make(chan bool)
	<-forever
}

func sendVote(timeout Timeout) {
	for {
		select {
		case <-timeout.ticker.C:
			myip := GetLocalIP()
			for _, s := range servers {
				if myip != s {
					dst := s + "_election"
					err := ch.Publish(
						dst,   // exchange
						dst,   // routing key
						false, // mandatory
						false, // immediate
						amqp.Publishing{
							ContentType: "text/plain",
							Body:        []byte(myip),
						})
					if err != nil {
						failOnError(err, "election publish error")
					}
				}
			}
			timeout.reset()
		}
	}
}

func listenWorkQ(workMsgs <-chan amqp.Delivery) {
	go func() {
		for d := range workMsgs {
			log.Printf("Received message: %s", d.Body)
		}
	}()

	fmt.Println("Service listening for events...")
}

func listenElectionQ(electionMsgs <-chan amqp.Delivery) {
	go func() {
		for d := range electionMsgs {
			log.Printf("Received message: %s", d.Body)
		}
	}()

	fmt.Println("Service listening for events...")
}

func listenHeartbeatQ(heartbeatMsgs <-chan amqp.Delivery) {
	go func() {
		for d := range heartbeatMsgs {
			log.Printf("Received message: %s", d.Body)
		}
	}()

	fmt.Println("Service listening for events...")
}
