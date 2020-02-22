package node

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"../spec"
)

var self spec.Self
var block = make(chan int, 1)

func Live(logf string) {
	// Initialize logging to file
	f, err := os.OpenFile(logf, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	// Get initial membership info
	self := spec.GetSelf()
	spec.ReportOnline(self.PID)

	go subscribeMembership()
	go listenForLeave()
	<-block
}

// Periodically poll for membership information
func subscribeMembership() {
	for {
		self = spec.GetSelf()
		time.Sleep(time.Second * spec.MembershipInterval)
	}
}

// Detect ctrl-c signal interrupts and dispatch [LEAVE]s to monitors accordingly
func listenForLeave() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(0)
	}()
}
