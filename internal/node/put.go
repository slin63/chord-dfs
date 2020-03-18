// Client and server stubs for RPCs.
package node

import (
	"errors"
	"fmt"
	"log"
	"time"

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
func Put(args *spec.PutArgs) []int {
	// Identify PID of server to give file to by calculating file's hash (FPID)
	FPID := hashing.MHash(args.Filename, self.M)
	PID := spec.NearestPID(FPID, &self)
	args.From = self.PID

	// Dispatch PutAssign RPC or perform on self
	if PID != self.PID {
		return []int{callPutAssign(PID, args)}
	}
	replicas, err := _putAssign(args)
	if err != nil {
		log.Fatal(err)
	}

	return replicas
}

// PutAssign (initiates PutAssign RPC on given PID)
func callPutAssign(PID int, args *spec.PutArgs) int {
	client := connect(PID)
	defer client.Close()

	var replicas []int
	if err := (*client).Call("Filesystem.PutAssign", *args, &replicas); err != nil {
		log.Fatal(err)
	}
	return PID
}

// PutAssign (receive RPC) (from: another server)
// Receive information about a file from another server
// Store that file on this machine and its replica nodes
// Return a slice of PIDs of servers with that file
func (f *Filesystem) PutAssign(args spec.PutArgs, replicas *[]int) error {
	var err error
	*replicas, err = _putAssign(&args)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

// Store the filename and data on this machine
// Also dispatch RPC calls to replica nodes
func _putAssign(args *spec.PutArgs) ([]int, error) {
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

	// Dispatch to replicas IF we are the main target for this file sharding
	if !args.Replicate {
		return []int{}, nil
	}

	replicaCh := make(chan int)
	replicas := []int{}
	replicasExpected := 0
	targetPID := self.PID
	spec.SelfRWMutex.RLock()
	for i := 0; i < config.C.Replicas; i++ {
		targetPID = spec.GetSuccessor(&self, targetPID)
		if targetPID != self.PID {
			go dispatchReplica(targetPID, args, replicaCh)
			replicasExpected++
		}
	}
	spec.SelfRWMutex.RUnlock()

	// Actually write to our own filesystem
	writes <- spec.WriteCmd{Name: args.Filename, Data: args.Data}

	// Watch for replicas coming in
	for {
		select {
		case replica := <-replicaCh:
			replicas = append(replicas, replica)
			if len(replicas) == replicasExpected {
				// Also include this node's self PID in the list of replicas.
				replicas = append(replicas, self.PID)
				return replicas, nil
			}
		case <-time.After(time.Duration(config.C.RPCTimeout) * 2 * time.Second):
			return replicas, errors.New("Timed out while waiting to hear back from all replica nodes.")
		}
	}
}

// Try and replicate the file inside args onto server PID.
func dispatchReplica(PID int, args *spec.PutArgs, resp chan<- int) {
	args.From = self.PID
	args.Replicate = false
	select {
	case resp <- callPutAssign(PID, args):
		config.LogIf(
			fmt.Sprintf("[FILEREPL] SUCCESSFULLY replicated file %s to [PID=%d]",
				args.Filename,
				PID,
			),
			config.C.LogFileReplication)
	case <-time.After(time.Duration(config.C.RPCTimeout) * time.Second):
		config.LogIf(
			fmt.Sprintf("[FILEREPL-X] FAILED to replicate file %s to [PID=%d]",
				args.Filename,
				PID,
			),
			config.C.LogFileReplication)
	}
}
