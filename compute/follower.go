package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func listenWorkQueue(server *Server) {
	conn, err := amqp.Dial("amqp://guest:guest@172.18.0.2:5672/")
	failOnError(err, "Error connecting to the broker")
	// Make sure we close the connection whenever the program is about to exit.
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	// Make sure we close the channel whenever the program is about to exit.
	defer ch.Close()

	// Subscribe to the queue.
	msgs, err := ch.Consume(
		"work_queue", // queue
		"",           // consumer id - empty means a random, unique id will be assigned
		false,        // auto acknowledgement of message delivery
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register as a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received message: %s", d.Body)

			close(forever)

			d.Ack(false)
		}
	}()

	fmt.Println("Service listening for events...")

	// Block until 'forever' receives a value, which will never happen.
	<-forever
}
