package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Service listening for HTTP...")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
