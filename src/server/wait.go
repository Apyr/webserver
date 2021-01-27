package server

import (
	"os"
	"os/signal"
)

func WaitInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
