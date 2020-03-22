package node

import (
	"errors"
	"fmt"
	"net/rpc"

	"github.com/slin63/chord-dfs/internal/config"
	"github.com/slin63/chord-dfs/internal/filesys"
	"github.com/slin63/chord-dfs/internal/hashing"
	"github.com/slin63/chord-dfs/internal/spec"
)

// Typical flow:
//   1. Client passes entry to Raft leader
//   2. Raft leader replicates entry and passes back to Filesystem.Execute RPC

// Delete (function to initiate DELETE action) (from: execute/client)
//   - Check if the file is stored on this node. If so, delete from self
//   - Hash the file onto some appropriate PID on the ring.
//   - Add that PID to a slice with PIDs of its N = Config.C.Replicas successors.
//   - Message those PIDs and tell them to delete their entries.
//   - If the connecting to the node fails:
//      - Add to retry queue # TODO: Implement

func Delete(args *spec.DeleteArgs) error {
	var client *rpc.Client
	var err error

	// Identify PID of server with file by calculating file's hash (FPID)
	FPID := hashing.MHash(args.Filename, self.M)
	PID := spec.NearestPID(FPID, &self)
	next := PID
	// Either:
	//   - Delete from own server
	//   - Delete from other server
	//   - Save this Delete action to be tried again later
	for i := 0; i < config.C.Replicas+1; i++ {
		if next == self.PID {
			// Just delete this file if we have it
			storeRWMutex.RLock()
			_, ok := store[args.Filename]
			if ok {
				delete(store, args.Filename)
				filesys.Remove(args.Filename)
			}
			storeRWMutex.RUnlock()
			next = spec.GetSuccessor(&self, next)
			continue
		}

		client, err = connectTimeout(next, config.C.RPCTimeout)
		if err != nil {
			config.LogIf(
				fmt.Sprintf("[DELETE-X] [PID=%d] Not responding. Skipping.", next),
				config.C.LogDelete,
			)
			next = spec.GetSuccessor(&self, next)
			continue
		}

		defer client.Close()
		err = callDelete(next, args, client)
		if err != nil {
			config.LogIf(
				fmt.Sprintf("[DELETE-X] [PID=%d] Error received from callDelete(). Skipping.", next),
				config.C.LogDelete,
			)
		}
		next = spec.GetSuccessor(&self, next)
	}

	return nil
}

// callDelete (initiates DeleteRespond RPC on given PID)
func callDelete(PID int, args *spec.DeleteArgs, client *rpc.Client) error {
	var b bool // RPC placeholder
	if err := (*client).Call("Filesystem.DeleteRespond", *args, &b); err != nil {
		return err
	}
	return nil
}

// DeleteRespond (receive RPC) (from: another server)
// Try and delete the given file
func (f *Filesystem) DeleteRespond(args spec.DeleteArgs, b *bool) error {
	storeRWMutex.RLock()
	_, ok := store[args.Filename]
	storeRWMutex.RUnlock()
	if ok {
		delete(store, args.Filename)
		return filesys.Remove(args.Filename)
	}

	return errors.New(fmt.Sprintf("File [%s] not found.", args.Filename))
}
