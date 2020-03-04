package main

import (
	"net"
)

func listenQueues(server *Server) {
	conn := Connect("amqp://guest:guest@172.18.0.2:5672/")
	// Make sure we close the connection whenever the program is about to exit.
	defer conn.Close()

	ch := GetChannel(conn)
	// Make sure we close the channel whenever the program is about to exit.
	defer ch.Close()

	ip := GetLocalIP()
	workQ := QueueDef{
		Exchange: ip,
		Queue:    ip + "_work",
		Binding:  ip,
	}
	DeclareExchange(ch, GetLocalIP())
	DeclareQueue(ch, workQ)

	// Subscribe to the queue.
	msgs := Subscribe(ch, workQ.Queue)
	ListenQueue(msgs)

	forever := make(chan bool)
	<-forever
}

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
