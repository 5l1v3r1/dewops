package main

import "log"

func failOnError(err error, msg string) {
	if err != nil {
		// Exit the program.
		log.Fatalf("%s: %s", msg, err)
	}
}
