package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

var ch *amqp.Channel

func listenResultQueue(server *Server) {
	conn := Connect("amqp://guest:guest@172.18.0.2:5672/")
	// Make sure we close the connection whenever the program is about to exit.
	defer conn.Close()

	ch = GetChannel(conn)
	// Make sure we close the channel whenever the program is about to exit.
	defer ch.Close()

	result := GetLocalIP() + "result"
	resultMsgs := CreateAndSubscribeQueue(ch, QueueDef{
		Exchange: result,
		Queue:    result,
		Binding:  result,
	}, "topic")
	listenResultQ(resultMsgs)

	sendHeartbeat(server.heartbeatTimeout)

	forever := make(chan bool)
	<-forever
}

func listenResultQ(resultMsgs <-chan amqp.Delivery) {
	go func() {
		for d := range resultMsgs {
			log.Printf("Received message: %s", d.Body)
		}
	}()

	fmt.Println("Service listening for events...")
}

func sendHeartbeat(timeout Timeout) {
	for {
		select {
		case <-timeout.ticker.C:
			myip := GetLocalIP()
			for _, s := range servers {
				if myip != s {
					dst := s + "_heart"
					err := ch.Publish(
						dst,   // exchange
						dst,   // routing key
						false, // mandatory
						false, // immediate
						amqp.Publishing{
							ContentType: "text/plain",
							Body:        []byte("leader"),
						})
					if err != nil {
						failOnError(err, "heartbeat publish error")
					}
				}
			}
			timeout.reset()
		}
	}
}
