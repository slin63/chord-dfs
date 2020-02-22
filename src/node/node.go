package node

import (
	"log"
	"os"

	"../spec"
)

var self spec.Self

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
	log.Println("hello")
	// So the program doesn't die
	// var wg sync.WaitGroup
	// wg.Add(1)

	// wg.Wait()
}
