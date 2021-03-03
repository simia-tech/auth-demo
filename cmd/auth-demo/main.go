package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/simia-tech/env"

	authdemo "github.com/simia-tech/auth-demo"
)

var (
	listenNetwork = env.String("LISTEN_NETWORK", "tcp")
	listenAddress = env.String("LISTEN_ADDRESS", "localhost:0")
)

func main() {
	s, err := authdemo.NewService(listenNetwork.Get(), listenAddress.Get())
	if err != nil {
		log.Fatal(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	if err := s.Close(); err != nil {
		log.Fatal(err)
	}
}
