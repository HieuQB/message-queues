package main

import (
	log "github.com/sirupsen/logrus"
	"ms-consumer"
	"os"
	"os/signal"
	"syscall"
)

var (
	name    = "consumer"
	version = "?"
)

func main() {
	stopCh := make(chan struct{})
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s := <-c
		log.Infof("signal received: %v. Shutting down service ...", s)
		close(stopCh)
	}()

	s := consumer.NewService(name, version)
	if err := s.Start(stopCh); err != nil {
		log.WithError(err).Error("failed to start service")
	}
}
