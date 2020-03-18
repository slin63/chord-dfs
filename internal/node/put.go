// Client and server stubs for RPCs.
package node

import (
	"fmt"
	"log"

	"github.com/slin63/chord-dfs/internal/config"
	"github.com/slin63/chord-dfs/internal/filesys"
	"github.com/slin63/chord-dfs/internal/hashing"
	"github.com/slin63/chord-dfs/internal/spec"
)

// Put (from: client)
//   - Hash the file onto some appropriate point on the ring.
//   - Message that point on the ring with the filename and data.
//   - Respond to the client with the process ID of the server that was selected.
func (f *Filesystem) Put(args spec.PutArgs, PIDPtr *int) error {
	if self.M != 0 {
		FPID := hashing.MHash(args.Filename, self.M)
		PIDPtr = spec.NearestPID(FPID, &self)

		// Dispatch PutAssign RPC or perform on self
		if *PIDPtr != self.PID {
			args.From = self.PID
			callPutAssign(*PIDPtr, &args)
		} else {
			_putAssign(&args)
		}
	}
	return nil
}

func callPutAssign(PID int, args *spec.PutArgs) {
	client := connect(PID)
	defer client.Close()

	var replicas []int
	if err := (*client).Call("Filesystem.PutAssign", *args, &replicas); err != nil {
		log.Fatal(err)
	}
}

// PutAssign (from: another server)
// Receive information about a file from another server
// Store that file on this machine and its replica nodes
// Return a slice of PIDs of servers with that file
func (f *Filesystem) PutAssign(args spec.PutArgs, replicas *[]int) error {
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
	filesys.Write(args.Filename, args.Data)

	return nil
}
