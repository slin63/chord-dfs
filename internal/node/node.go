// Stuff the server needs to do to stay alive and do its job.
package node

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/slin63/chord-dfs/internal/config"
	"github.com/slin63/chord-dfs/internal/spec"
)

var self spec.Self
var block = make(chan int, 1)

// Maps filename to length of byte array
var store = make(map[string]int)

func Live() {
	// Create directory for storing files
	os.Mkdir(config.C.Filedir, 0644)

	// Get initial membership info
	spec.GetSelf(&self)
	log.SetPrefix(config.C.Prefix + fmt.Sprintf(" [PID=%d]", self.PID) + " - ")
	spec.ReportOnline()

	go serveFilesystemRPC()
	go subscribeMembership()
	go listenForLeave()
	<-block
}

// Periodically poll for membership information
func subscribeMembership() {
	for {
		spec.GetSelf(&self)
		time.Sleep(time.Second * time.Duration(config.C.MemberInterval))
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