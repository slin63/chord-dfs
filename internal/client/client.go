package client

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/rpc"
	"strconv"
	"strings"

	"github.com/slin63/chord-dfs/internal/config"
	"github.com/slin63/chord-dfs/pkg/parser"

	"github.com/slin63/raft-consensus/pkg/responses"
)

const HelpS = `Available operations:
1. put <localfilename> <sdfsfilename> (from local dir)
2. get <sdfsfilename> <localfilename> (fetches to local dir)
3. delete <sdfsfilename>
4. ls <filename> (list all machines where this data is stored)
5. store <PID> (list all files stored on this machine)`

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
	var sdfs string
	if len(args) == 0 {
		fmt.Println(HelpS)
		return
	}

	// Check input validity. If valid, send off to Raft for replication.
	method, _, ok := parser.ParseEntry(args)
	if !ok {
		fmt.Println("Invalid input!")
		fmt.Println(HelpS)
		return
	}

	// Try connecting to any node.
	client, err := getLiveNode()
	if err != nil {
		fmt.Println(err)
		return
	}

	switch method {
	case parser.PUT:
		local = args[1]
		f, err := ioutil.ReadFile(local)
		if err != nil {
			fmt.Println(err)
			return
		}
		args = append(args, string(f))
		methodS, _ = parser.MethodString(parser.PUT)
	case parser.GET:
		sdfs = args[1]
		local = args[2]
		methodS, _ = parser.MethodString(parser.GET)
	case parser.DELETE:
		sdfs = args[1]
		methodS, _ = parser.MethodString(parser.DELETE)
	case parser.LS:
		sdfs = args[1]
		methodS, _ = parser.MethodString(parser.LS)
	case parser.STORE:
		_, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		methodS, _ = parser.MethodString(parser.STORE)
	default:
		panic("TODO: Add more methods")
	}

	// Try sending a PutEntry.
	result, err := callPutEntry(client, args)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !result.Success && result.Error != responses.LEADERREDIRECT {
		fmt.Printf(
			"Something went wrong inside the massive black box. Sorry!\n\terrcode: %d\n\tdata: %s\n",
			result.Error, result.Data,
		)
		return
	}

	resultsFormattedS := "Result: %s\n\tEntry: %s"
	resultsFormatted := fmt.Sprintf(resultsFormattedS, formatEntry(result.Data), formatEntry(result.Entry))
	log.Printf("Success!\n\tExecuted %s.\n\t%s", methodS, resultsFormatted)

	// Write to local file system
	if method == parser.GET {
		err := ioutil.WriteFile(local, []byte(result.Data), 0644)
		if err != nil {
			fmt.Printf("Error trying to write file to local filesystem: %v", err)
			return
		}

		log.Printf("Wrote [SDFS=%s] to [LOCAL=%s].", sdfs, local)
	}
}

func formatEntry(s string) string {
	s = strings.ReplaceAll(s, "\n", "; ")
	if len(s) > 40 {
		s = s[0:40] + " ... " + fmt.Sprintf("(+%d)", len(s)-40)
	}
	return s
}

// Send a PutEntry to the given server. Possible outcomes:
//  - Server is the leader and PutEntry is accepted
//  - Server is not the leader and offers a redirect to the leader
//    - Connect to that server and set off the PutEntry
func callPutEntry(client *rpc.Client, args []string) (*responses.Result, error) {
	var result responses.Result
	var err error
	// Dispatch to the given client
	if err = client.Call("Ocean.PutEntry", strings.Join(args, " "), &result); err != nil {
		fmt.Println(err)
		return &result, err
	}

	// Redispatch request to proper leader. Assume we're never redirected more than once
	if !result.Success && result.Error == responses.LEADERREDIRECT {
		redir := strings.Split(result.Data, ",")
		addr := redir[1]
		fmt.Printf("[CALLPUTENTRY] Received leader redirect: %s\n", addr)
		client, err = rpc.DialHTTP("tcp", addr+":"+config.C.RaftRPCPort)
		if err != nil {
			fmt.Println("[ERROR] PutEntry() dialing:\n", err)
			return &result, err
		}

		if err = client.Call("Ocean.PutEntry", strings.Join(args, " "), &result); err != nil {
			fmt.Println(err)
			return &result, err
		}
		fmt.Printf("[CALLPUTENTRY] Successfully dialed: %s\n", addr)
	}

	return &result, nil
}

// Return any node that is alive, trying the introducer first.
func getLiveNode() (*rpc.Client, error) {
	client, err := rpc.DialHTTP("tcp", config.C.Introducer+":"+config.C.RaftRPCPort)
	if err != nil {
		fmt.Println("[GETLIVENODE] Introducer dead. Trying to find new leader")
		for i := 0; i < config.C.MaxServerLookups; i++ {
			next := fmt.Sprintf("%s%d", config.C.ServerPrefix, i+1)
			client, err = rpc.DialHTTP("tcp", next+":"+config.C.RaftRPCPort)
			if err == nil {
				fmt.Printf("[GETLIVENODE] Alive node at %s\n", next)
				break
			}
			fmt.Printf("[GETLIVENODE] %s not available.\n", next)
		}
	}

	return client, err
}
