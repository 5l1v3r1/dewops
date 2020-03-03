package main

func listenWorkQueue(server *Server) {
	conn := connect("amqp://guest:guest@172.18.0.2:5672/")
	// Make sure we close the connection whenever the program is about to exit.
	defer conn.Close()

	ch := getChannel(conn)
	// Make sure we close the channel whenever the program is about to exit.
	defer ch.Close()

	// Subscribe to the queue.
	msgs := subscribe(ch, "work_queue")

	listenQueue(msgs)
}
