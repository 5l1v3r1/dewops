package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"strings"
)

var leaderIP string

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
	listenWorkQ(ch, workMsgs)

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
	listenHeartbeatQ(hbMsgs, server)

	//sendVote(ch, server)

	forever := make(chan bool)
	<-forever
}

func sendVote(ch *amqp.Channel, server *Server) {
	timeout := server.electionTimeout
	myip := GetLocalIP()
	server.term++
	body := myip + " " + string(server.term)
	for {
		select {
		case <-timeout.ticker.C:
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
							Body:        []byte(body),
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

func listenWorkQ(ch *amqp.Channel, workMsgs <-chan amqp.Delivery) {
	go func() {
		for d := range workMsgs {
			log.Printf("Received message: %s", d.Body)
			sendResult(ch, wordCount(string(d.Body)))
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

func listenHeartbeatQ(heartbeatMsgs <-chan amqp.Delivery, server *Server) {
	go func() {
		for d := range heartbeatMsgs {
			select {
			case <-server.heartbeatTimeout.ticker.C:
				leaderIP = string(d.Body)
				server.electionTimeout.reset()
				log.Printf("Received message: %s", d.Body)
				server.heartbeatTimeout.reset()
			}
		}
	}()

	fmt.Println("Service listening for events...")
}

func wordCount(words string) map[string]int {
	wordList := strings.Fields(words)
	wordFreq := make(map[string]int)

	for _, word := range wordList {
		_, ok := wordFreq[word]

		if ok == true {
			wordFreq[word] += 1
		} else {
			wordFreq[word] = 1
		}
	}

	return wordFreq
}

func sendResult(ch *amqp.Channel, res map[string]int) {
	resQ := leaderIP + "_result"
	jsonString, err := json.Marshal(res)
	if err != nil {
		failOnError(err, "map to json failed")
	}
	err = ch.Publish(
		resQ,   // exchange
		resQ,   // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			MessageId:   myIP,
			Body:        jsonString,
		})
	if err != nil {
		failOnError(err, "election publish error")
	}
}
