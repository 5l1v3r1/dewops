package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"strings"
)

func connect(url string) *amqp.Connection {
	connection, err := amqp.Dial(url)
	failOnError(err, "Error connecting to the broker")
	return connection
}

func getChannel(connection *amqp.Connection) *amqp.Channel {
	channel, err := connection.Channel()
	failOnError(err, "Failed to open a channel")
	return channel
}

func subscribe(ch *amqp.Channel, qname string) <-chan amqp.Delivery {
	messages, err := ch.Consume(
		qname, // queue
		"",    // consumer id - empty means a random, unique id will be assigned
		false, // auto acknowledgement of message delivery
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register as a consumer")

	return messages
}

func listenQueue(msgs <-chan amqp.Delivery) {
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received message: %s", d.Body)
			if strings.Contains(string(d.Body), "stop") {
				close(forever)
			}

			d.Ack(false)
		}
	}()

	fmt.Println("Service listening for events...")

	// Block until 'forever' receives a value, which will never happen.
	<-forever
}
