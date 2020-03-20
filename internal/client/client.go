package client

import (
	"fmt"
	"io/ioutil"
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
	var methodS string
	var local string
	if len(args) == 0 {
		fmt.Println(helpS)
		return
	}

	// Check input validity. If valid, send off to Raft for replication.
	method, _, ok := parser.ParseEntry(args)
	if !ok {
		fmt.Println("Invalid input!")
		log.Fatal(helpS)
		return
	}

	client, err := rpc.DialHTTP("tcp", "localhost:"+config.C.RaftRPCPort)
	if err != nil {
		log.Fatal("[ERROR] PutEntry() dialing:", err)
	}

	switch method {
	case parser.PUT:
		local = args[1]
		f, err := ioutil.ReadFile(local)
		if err != nil {
			log.Fatal("[putArgs()]: ", err)
		}
		args = append(args, string(f))
		methodS, _ = parser.MethodString(parser.PUT)
	default:
		panic("TODO: Add more methods")
	}

	var result *responses.Result
	if err = client.Call("Ocean.PutEntry", strings.Join(args, " "), &result); err != nil {
		log.Fatal(err)
	}

	if !result.Success {
		log.Println("Something went wrong inside the massive black box. Sorry!")
		return
	}

	log.Printf("Success!\n\tExecuted %s on %s.\n\tResult: [%s]", methodS, local, result.Data)
}
