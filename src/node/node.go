package node

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"
	"time"

	"../hashing"
	"../spec"
)

var self spec.Self
var block = make(chan int, 1)

// RPC type
type Filesystem int

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

	go serveFilesystemRPC()
	go subscribeMembership()
	go listenForLeave()
	<-block
}

func serveFilesystemRPC() {
	fs := new(Filesystem)
	rpc.Register(fs)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":"+spec.FilesystemRPCPort)
	if e != nil {
		log.Fatal("[ERROR] serveFilesystemRPC():", e)
	}
	go http.Serve(l, nil)
}

// Hash the file onto some appropriate point on the ring.
// Respond to the client with the PID of the server that was selected.
func (f *Filesystem) Put(args spec.PutArgs, PID *int) error {
	log.SetPrefix("Put(): ")
	defer log.SetPrefix(spec.Prefix)
	if self.M != 0 {
		FPID := hashing.GetPID(args.Filename, self.M)
		log.Println("FPID: ", FPID)

		// success, handle errors here
		PID = spec.GetSuccPID(FPID, &self)
		log.Println("PID: ", *PID)
	}
	return nil
}

// Periodically poll for membership information
func subscribeMembership() {
	for {
		spec.SelfSem.Lock()
		self = spec.GetSelf()
		spec.SelfSem.Unlock()

		time.Sleep(time.Second * spec.MemberInterval)
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
