package main

func LBCommunication(server *Server) {

}

func listenResultQueue() {
	conn := Connect("amqp://guest:guest@172.18.0.2:5672/")
	// Make sure we close the connection whenever the program is about to exit.
	defer conn.Close()

	ch := GetChannel(conn)
	// Make sure we close the channel whenever the program is about to exit.
	defer ch.Close()

	result := GetLocalIP() + "result"
	CreateAndSubscribeQueue(ch, QueueDef{
		Exchange: result,
		Queue:    result,
		Binding:  result,
	}, RESULTQ, "topic")

	forever := make(chan bool)
	<-forever
}
