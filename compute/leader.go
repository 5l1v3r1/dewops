package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"math"
)

var myIP = GetLocalIP()

func listenResultQueue(server *Server) {
	conn := Connect("amqp://guest:guest@172.18.0.2:5672/")
	// Make sure we close the connection whenever the program is about to exit.
	defer conn.Close()

	ch := GetChannel(conn)
	// Make sure we close the channel whenever the program is about to exit.
	defer ch.Close()

	result := GetLocalIP() + "_result"
	resultMsgs := CreateAndSubscribeQueue(ch, QueueDef{
		Exchange: result,
		Queue:    result,
		Binding:  result,
	}, "topic")
	listenResultQ(resultMsgs)

	sendWork(ch)

	sendHeartbeat(ch, server.heartbeatTimeout)
	forever := make(chan bool)
	<-forever
}

func listenResultQ(resultMsgs <-chan amqp.Delivery) {
	go func() {
		for d := range resultMsgs {
			log.Printf("%s Received message: %s", d.MessageId, d.Body)
		}
	}()

	fmt.Println("Service listening for events...")
}

func sendHeartbeat(ch *amqp.Channel, timeout Timeout) {
	for {
		select {
		case <-timeout.ticker.C:
			for _, s := range servers {
				if myIP != s {
					dst := s + "_heart"
					sendMessage(ch, dst, myIP)
				}
			}
			timeout.reset()
		}
	}
}

func sendWork(ch *amqp.Channel) {
	log.Printf("Work sending is started...")
	data, err := ioutil.ReadFile("alice.txt")
	//log.Printf(string(data))
	if err != nil {
		failOnError(err, "read file error")
	}

	datas := string(data)

	splittedWork := splitByWidthMake(datas, len(datas) / 4)

	i := 0
	for _, s := range servers {
		if myIP != s {
			dst := s + "_work"
			sendMessage(ch, dst, splittedWork[i])
			i++
		}
	}
}

func sendMessage(ch *amqp.Channel, dst string, data string) {
	err := ch.Publish(
		dst,   // exchange
		dst,   // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(data),
			})
	if err != nil {
		failOnError(err, "send message publish error")
	}
}

func splitByWidthMake(str string, size int) []string {
	strLength := len(str)
	splitedLength := int(math.Ceil(float64(strLength) / float64(size)))
	splited := make([]string, splitedLength)
	var start, stop int
	for i := 0; i < splitedLength; i += 1 {
		start = i * size
		stop = start + size
		if stop > strLength {
			stop = strLength
		}
		splited[i] = str[start : stop]
	}
	return splited
}