// Client and server stubs for RPCs.
package node

import (
	"fmt"
	"log"

	"github.com/slin63/chord-dfs/internal/config"
	"github.com/slin63/chord-dfs/internal/hashing"
	"github.com/slin63/chord-dfs/internal/spec"
)

// Typical flow:
//   1. Client passes entry to Raft leader
//   2. Raft leader replicates entry and passes back to Filesystem.Execute RPC

// From there:
// -> Filesystem.Execute() calls Put(File, PID)
//   -> Server receives Put(File, PID)
//      -> PID == self.PID
//         -> Server calls _putAssign(File, PID) on self
//            -> Return result to Execute, Execute returns result to Raft leader
//               -> Raft leader returns result to Client 🥳
//      -> PID != self.PID
//         -> Server calls callPutAssign(File, PID) on ServerB with PID == PID
//            -> ServerB receives PutAssign(File, PID)
//               -> ServerB calls _putAssign(File, PID) on self
//                 -> Return result to Execute, Execute returns result to Raft leader
//                    -> Raft leader returns result to Client 🥳

// Put (function to initiate PUT action) (from: client)
//   - Hash the file onto some appropriate point on the ring.
//   - Message that point on the ring with the filename and data.
//   - Respond to the client with the process ID of the server that was selected.
func Put(args *spec.PutArgs) error {
	if self.M != 0 {
		// Identify PID of server to give file to by calculating file's hash (FPID)
		FPID := hashing.MHash(args.Filename, self.M)
		PID := spec.NearestPID(FPID, &self)

		// Dispatch PutAssign RPC or perform on self
		if PID != self.PID {
			args.From = self.PID
			callPutAssign(PID, args)
		} else {
			_putAssign(args)
		}
	}
	return nil
}

// PutAssign (initiates PutAssign RPC on given PID)
func callPutAssign(PID int, args *spec.PutArgs) {
	client := connect(PID)
	defer client.Close()

	var replicas []int
	if err := (*client).Call("Filesystem.PutAssign", *args, &replicas); err != nil {
		log.Fatal(err)
	}
}

// PutAssign (receive RPC) (from: another server)
// Receive information about a file from another server
// Store that file on this machine and its replica nodes
// Return a slice of PIDs of servers with that file
func (f *Filesystem) PutAssign(args spec.PutArgs, replicas *[]int) error {
	_putAssign(&args)
	// TODO (03/18 @ 15:41): find a way to return nodes with replicas
	return nil
}

// Store the filename and data on this machine
// Also dispatch RPC calls to replica nodes
func _putAssign(args *spec.PutArgs) {
	storeRWMutex.Lock()
	defer storeRWMutex.Unlock()

	v, ok := store[args.Filename]
	if ok {
		config.LogIf(
			fmt.Sprintf("[PUT] Updating %s:%d -> %d", args.Filename, v, len(args.Data)), config.C.LogPutAssign,
		)
	} else {
		config.LogIf(
			fmt.Sprintf("[PUT] Setting %s:%d", args.Filename, len(args.Data)), config.C.LogPutAssign,
		)
	}
	// Update in memory store
	store[args.Filename] = len(args.Data)

	// Actually write to filesystem
	writes <- spec.WriteCmd{Name: args.Filename, Data: args.Data}
}
