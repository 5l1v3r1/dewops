package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"net"
	"strings"
)

type QueueName uint

const (
	WORKERQ QueueName = iota
	ELECTIONQ
	RESULTQ
	LEADERQ
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

func DeclareExchange(channel *amqp.Channel, exchangeName string, exchangeType string) {
	// Create the exchange if it doesn't already exist.
	err := channel.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type
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

func ListenQueue(msgs <-chan amqp.Delivery, qn QueueName) {
	go func() {
		for d := range msgs {
			log.Printf("Received message: %s", d.Body)
			if strings.Contains(string(d.Body), "stop") {
				break
			}
			if qn == WORKERQ {
				log.Printf("This is workerq")
			} else if qn == ELECTIONQ {
				log.Printf("This is electionq")
			} else if qn == RESULTQ {
				log.Printf("This is resultq")
			} else if qn == LEADERQ {
				log.Printf("This is leaderq")
			}

			d.Ack(false)
		}
	}()

	fmt.Println("Service listening for events...")
}

func CreateAndSubscribeQueue(ch *amqp.Channel, def QueueDef, qn QueueName, exchangeType string) {
	DeclareExchange(ch, def.Exchange, exchangeType)
	DeclareQueue(ch, def)

	// Subscribe to the queue.
	msgs := Subscribe(ch, def.Queue)
	ListenQueue(msgs, qn)
}

// local ip is used for generating unique queue names for each service
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func failOnError(err error, msg string) {
	if err != nil {
		// Exit the program.
		log.Fatalf("%s: %s", msg, err)
	}
}
