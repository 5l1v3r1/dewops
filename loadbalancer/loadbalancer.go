package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello")
		listenLeaderQueue()
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func listenLeaderQueue() {
	conn := Connect("amqp://guest:guest@172.18.0.2:5672/")
	// Make sure we close the connection whenever the program is about to exit.
	defer conn.Close()

	ch := GetChannel(conn)
	// Make sure we close the channel whenever the program is about to exit.
	defer ch.Close()

	// Subscribe to the queue.
	leaderMsgs := CreateAndSubscribeQueue(ch, QueueDef{
		Exchange: "leaderq",
		Queue:    "leaderq",
		Binding:  "leaderq",
	}, "topic")

	go func() {
		for d := range leaderMsgs {
			log.Printf("Received message: %s", d.Body)
		}
	}()

	fmt.Println("Service listening for events...")
}

/*https://softwareengineering.stackexchange.com/questions/312956/what-does-a-load-balancer-return
The end-IP is not published.
The process actually works in a way the client (a user hitting the balancer)
believes they are communicating with the balancer, while talking to an actual node.

In a very simple explanation, most transactions work like this:

A user makes request to the load balancer.
The balancer decides which node is the most suitable (based on the strategy you're using for balancing) and choses (changes) the destination IP.
(This is where the magic happens.) The node receives a request, accepts the connection and responds back to the balancer.
The balancer changes the response IP back to a virtual one, the one of the balancer and forwards the response to the user.
Voilà, the user receives response with the IP of the initial request, even though it was actually processed somewhere else.
Keep in mind the packet rewriting (the change of the IP address in the step 4) is very important. Without it the client, receiving a packet from an IP it does not trust, would simply discard the response.
*/
