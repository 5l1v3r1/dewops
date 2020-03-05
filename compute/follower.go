package main

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
	CreateAndSubscribeQueue(ch, QueueDef{
		Exchange: work,
		Queue:    work,
		Binding:  work,
	}, WORKERQ, "topic")

	// election queue
	elect := ip + "_election"
	CreateAndSubscribeQueue(ch, QueueDef{
		Exchange: elect,
		Queue:    elect,
		Binding:  elect,
	}, ELECTIONQ, "fanout")

	forever := make(chan bool)
	<-forever
}
