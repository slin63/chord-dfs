package client

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/rpc"

	"github.com/slin63/chord-dfs/internal/config"
	"github.com/slin63/chord-dfs/internal/spec"

	"github.com/slin63/raft-consensus/pkg/responses"
)

// Interfaces with Raft leader.
//   - x Client validates syntax of user entry
//   - x Sends entry to Raft leader
//   - x Raft leader replicates entry to replica nodes
//   - After successful replication, Raft leader tries applying change by contacting
//       DFS server via "handleEntry" RPC
//   - DFS server returns results to Raft leader,
//   - Raft leader returns results to Client.
func Parse(args []string) {
	if len(args) == 0 {
		fmt.Println(helpS)
		return
	}

	// Check input validity. If valid, send off to Raft for replication.
	if entry, ok := parseEntry(args); !ok {
		fmt.Println("Invalid input!")
		fmt.Println(helpS)
		return
	} else {
		client, err := rpc.DialHTTP("tcp", "localhost:"+config.C.RaftRPCPort)
		if err != nil {
			log.Fatal("[ERROR] PutEntry() dialing:", err)
		}

		// PID of assigned server
		var result *responses.Result
		if err = client.Call("Ocean.PutEntry", entry, &result); err != nil {
			log.Fatal(err)
		}
		log.Println(*result)
	}
}

// - Read file `local`
// - Send RPC to DFS server with correct arguments
// - DFS server decides what to do with it and where to put it
func put(local, sdfs string) {
	f, err := ioutil.ReadFile(local)
	if err != nil {
		log.Fatal("put(): ", err)
	}
	log.Println(f)

	client, err := rpc.DialHTTP("tcp", "localhost:"+config.C.FilesystemRPCPort)
	if err != nil {
		log.Fatal("put() dialing:", err)
	}

	// PID of assigned server
	var assigned int
	args := spec.PutArgs{sdfs, f, spec.NILPID}
	if err = client.Call("Filesystem.Put", args, &assigned); err != nil {
		log.Println(assigned)
	}
}
