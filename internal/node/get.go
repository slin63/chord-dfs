// Client and server stubs for RPCs.
package node

import (
	"errors"
	"fmt"
	"log"

	"github.com/slin63/chord-dfs/internal/filesys"
	"github.com/slin63/chord-dfs/internal/hashing"
	"github.com/slin63/chord-dfs/internal/spec"
)

// Typical flow:
//   1. Client passes entry to Raft leader
//   2. Raft leader replicates entry and passes back to Filesystem.Execute RPC

// Get (function to initiate PUT action) (from: execute/client)
//   - Check if the file is stored on this machine. If so, just return it.
//   - Hash the file onto some appropriate PID on the ring.
//   - Message that PID on the ring  with the filename.
//   - If it's alive:
//     - Get the file data from that PID and return it to the client
//   - If it's dead:
//     - Try checking its two successors until one of them responds
func Get(args *spec.GetArgs) ([]byte, error) {
	// Just return this file if we have it
	spec.SelfRWMutex.RLock()
	_, ok := store[args.Filename]
	spec.SelfRWMutex.RUnlock()
	if ok {
		return filesys.Read(args.Filename), nil
	}

	// Identify PID of server to get file to by calculating file's hash (FPID)
	FPID := hashing.MHash(args.Filename, self.M)
	PID := spec.NearestPID(FPID, &self)

	// Dispatch PutAssign RPC or perform on self
	// TODO (03/20 @ 16:50): Implement trying successors if this times out
	return callGet(PID, args)
}

// callGet (initiates GetRespond RPC on given PID)
func callGet(PID int, args *spec.GetArgs) ([]byte, error) {
	client := connect(PID)
	defer client.Close()

	var data []byte
	if err := (*client).Call("Filesystem.GetRespond", *args, &data); err != nil {
		log.Printf("[CALLGET-X] Error while retrieving file from [PID=%d]: %v", PID, err)
		return data, err
	}

	return data, nil
}

// GetRespond (receive RPC) (from: another server)
// Receive information about a file from another server
// Store that file on this machine and its replica nodes
// Return a slice of PIDs of servers with that file
func (f *Filesystem) GetRespond(args spec.GetArgs, data *[]byte) error {
	// Just return this file if we have it
	spec.SelfRWMutex.RLock()
	_, ok := store[args.Filename]
	spec.SelfRWMutex.RUnlock()
	if ok {
		*data = filesys.Read(args.Filename)
		return nil
	}

	return errors.New(fmt.Sprintf("File [%s] not found.", args.Filename))
}
