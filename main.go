package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"go.uber.org/automaxprocs/maxprocs"
)

var build = "develop"

func main() {
	if _, err := maxprocs.Set(); err != nil {
		fmt.Println("maxprocs error:", err)
		os.Exit(1)
	}

	g := runtime.GOMAXPROCS(0)

	log.Printf("Started shrt-api %s service! CPUs: %d\n", build, g)
	defer log.Printf("Stopped shrt-api %s service!\n", build)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	log.Println("Stopping shrt-api service...")
}
