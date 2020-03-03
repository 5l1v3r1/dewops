package main

import (
	"fmt"
	"log"
	"net/http"
)

func clientCommunication(server *Server) {
	fmt.Println("Service listening for HTTP...")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello")
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}
