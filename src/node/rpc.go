package node

import (
	"log"
	"net"
	"net/http"
	"net/rpc"

	"../hashing"

	"../spec"
)

type Filesystem int

func serveFilesystemRPC() {
	fs := new(Filesystem)
	rpc.Register(fs)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":"+spec.FilesystemRPCPort)
	if e != nil {
		log.Fatal("[ERROR] serveFilesystemRPC():", e)
	}
	http.Serve(l, nil)
}

// Put (from: client)
// Hash the file onto some appropriate point on the ring.
// Message that point on the ring with the filename and data.
// Respond to the client with the process ID of the server that was selected.
func (f *Filesystem) Put(args spec.PutArgs, PID *int) error {
	log.SetPrefix(spec.Prefix + "Put(): ")
	defer log.SetPrefix(spec.Prefix)
	if self.M != 0 {
		FPID := hashing.MHash(args.Filename, self.M)
		log.Println("FPID: ", FPID)

		// success, handle errors here
		PID = spec.GetSuccPID(FPID, &self)
		log.Println("PID: ", *PID)

		//
		// TODO (02/25 @ 14:14): handle case where pid  =  selfPid
		args.From = self.PID
		putAssignC(*PID, &args)
	}
	return nil
}

// PutAssign (from: server)
// Receive information about a file from another server
// Store that file on this machine and its replica nodes
// Return a slice of PIDs of servers with that file
func (f *Filesystem) PutAssign(args spec.PutArgs, replicas *[]int) error {
	log.SetPrefix(spec.Prefix + "PutAssign(): ")
	defer log.SetPrefix(spec.Prefix)
	log.Println(args)
	// TODO (02/25 @ 13:21): implement
	return nil
}

// - load in with read file
// - send over tcp (have to use tcp because we're using rpcs) to server
// - server decides what to do with it and where to put it
func putAssignC(PID int, args *spec.PutArgs) {
	log.SetPrefix(spec.Prefix + "putAssignC(): ")
	defer log.SetPrefix(spec.Prefix)
	client := connect(PID)
	defer client.Close()

	var replicas []int
	if err := (*client).Call("Filesystem.PutAssign", *args, &replicas); err != nil {
		log.Fatal(err)
	}
}

// Connect to some RPC server and return a pointer to the client
func connect(PID int) *rpc.Client {
	node, ok := self.MemberMap[PID]
	if !ok {
		log.Fatalf("[PID=%d] member not found.", PID)
	}
	client, err := rpc.DialHTTP("tcp", (*node).IP+":"+spec.FilesystemRPCPort)
	if err != nil {
		log.Fatal("put() dialing:", err)
	}

	return client
}
