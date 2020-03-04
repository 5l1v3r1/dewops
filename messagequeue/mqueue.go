package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"strings"
)

type QueueDef struct {
	Exchange string
	Queue    string
	Binding  string
}

func Connect(url string) *amqp.Connection {
	connection, err := amqp.Dial(url)
	failOnError(err, "Error connecting to the broker")
	return connection
}

func GetChannel(connection *amqp.Connection) *amqp.Channel {
	channel, err := connection.Channel()
	failOnError(err, "Failed to open a channel")
	return channel
}

func DeclareExchange(channel *amqp.Channel, exchangeName string) {
	// Create the exchange if it doesn't already exist.
	err := channel.ExchangeDeclare(
		exchangeName, // name
		"topic",      // type
		true,         // durable
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Error creating the exchange "+exchangeName)
}

func DeclareQueue(channel *amqp.Channel, def QueueDef) {
	// Create the queue if it doesn't already exist.
	// This does not need to be done in the publisher because the
	// queue is only relevant to the consumer, which subscribes to it.
	// Like the exchange, let's make it durable (saved to disk) too.
	q, err := channel.QueueDeclare(
		def.Queue, // name - empty means a random, unique name will be assigned
		true,      // durable
		false,     // delete when the last consumer unsubscribe
		false,
		false,
		nil,
	)
	failOnError(err, "Error creating the queue")

	// Bind the queue to the exchange based on a string pattern (binding key).
	err = channel.QueueBind(
		q.Name,       // queue name
		def.Binding,  // binding key
		def.Exchange, // exchange
		false,
		nil,
	)

	failOnError(err, "Error binding the queue")
}

func Subscribe(ch *amqp.Channel, qname string) <-chan amqp.Delivery {
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

func ListenQueue(msgs <-chan amqp.Delivery) {
	go func() {
		for d := range msgs {
			log.Printf("Received message: %s", d.Body)
			if strings.Contains(string(d.Body), "stop") {
				break
			}

			d.Ack(false)
		}
	}()

	fmt.Println("Service listening for events...")
}

func failOnError(err error, msg string) {
	if err != nil {
		// Exit the program.
		log.Fatalf("%s: %s", msg, err)
	}
}
