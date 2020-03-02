package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		// Exit the program.
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
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

			// Update the user's data on the service's
			// associated datastore using a local transaction...

			// The 'false' indicates the success of a single delivery, 'true' would
			// mean that this delivery and all prior unacknowledged deliveries on this
			// channel will be acknowledged, which I find no reason for in this example.
			d.Ack(false)
		}
	}()

	fmt.Println("Service listening for events...")

	// Block until 'forever' receives a value, which will never happen.
	<-forever
}
