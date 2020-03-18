package client

import (
	"fmt"
	"log"
	"net/rpc"
	"strings"

	"github.com/slin63/chord-dfs/internal/config"
	"github.com/slin63/chord-dfs/pkg/parser"

	"github.com/slin63/raft-consensus/pkg/responses"
)

const helpS = `Available operations:
1. put localfilename sdfsfilename (from local dir)
2. get sdfsfilename localfilename (fetches to local dir)
3. delete sdfsfilename
4. ls filename (list all machines where this data is stored)
5. store (list all files stored on this machine)`

// Interfaces with Raft leader.
//   - Client validates syntax of user entry
//   - Sends entry to Raft leader
//   - Raft leader replicates entry to replica nodes
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
	if _, _, ok := parser.ParseEntry(args); !ok {
		fmt.Println("Invalid input!")
		log.Fatal(helpS)
		return
	} else {
		client, err := rpc.DialHTTP("tcp", "localhost:"+config.C.RaftRPCPort)
		if err != nil {
			log.Fatal("[ERROR] PutEntry() dialing:", err)
		}

		// PID of assigned server
		var result *responses.Result
		if err = client.Call("Ocean.PutEntry", strings.Join(args, " "), &result); err != nil {
			log.Fatal(err)
		}
		log.Println(*result)
	}
}
