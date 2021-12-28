package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

var build = "develop"

func main() {
	log.Printf("Started shrt %s service!\n", build)
	defer log.Printf("Stopped shrt %s service!\n", build)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	log.Println("Stopping shrt service...")
}
