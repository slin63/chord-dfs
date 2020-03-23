package node

import (
	"errors"
	"fmt"
	"net/rpc"

	"github.com/slin63/chord-dfs/internal/config"
	"github.com/slin63/chord-dfs/internal/hashing"
	"github.com/slin63/chord-dfs/internal/spec"
)

// Typical flow:
//   1. Client passes entry to Raft leader
//   2. Raft leader replicates entry and passes back to Filesystem.Execute RPC

// Ls (function to initiate LS action) (from: execute/client)
//   - Hash the file onto some appropriate PID on the ring.
//   - Have that PID check itself and its replicas for the file.
//     - Return a slice with all the PIDs that have that file.
func Ls(args *spec.LsArgs) ([]int, error) {
	var client *rpc.Client
	var err error
	holders := []int{}

	// Identify PID of server with file by calculating file's hash (FPID)
	FPID := hashing.MHash(args.Filename, self.M)
	PID := spec.NearestPID(FPID, &self)
	next := PID
	// Either:
	//   - Ls from own server
	//   - Ls from other server
	for i := 0; i < config.C.Replicas+1; i++ {
		if next == self.PID {
			// Check if we have this file
			storeRWMutex.RLock()
			_, ok := store[args.Filename]
			if ok {
				holders = append(holders, self.PID)
			}
			storeRWMutex.RUnlock()
			next = spec.GetSuccessor(&self, next)
			continue
		}

		client, err = connectTimeout(next, config.C.RPCTimeout)
		if err != nil {
			config.LogIf(
				fmt.Sprintf("[LS-X] [PID=%d] Not responding. Skipping.", next),
				config.C.LogLs,
			)
			next = spec.GetSuccessor(&self, next)
			continue
		}

		defer client.Close()
		err = callLs(next, args, client)
		if err != nil {
			config.LogIf(
				fmt.Sprintf("[LS-X] [PID=%d] Error received from callLs(). Skipping.", next),
				config.C.LogLs,
			)
			next = spec.GetSuccessor(&self, next)
			continue
		}

		holders = append(holders, next)
		next = spec.GetSuccessor(&self, next)
	}

	return holders, nil
}

// callLs (initiates LsRespond RPC on given PID)
func callLs(PID int, args *spec.LsArgs, client *rpc.Client) error {
	var success bool
	if err := (*client).Call("Filesystem.LsRespond", *args, &success); err != nil {
		return err
	}
	return nil
}

// LsRespond (receive RPC) (from: another server)
// Try and delete the given file
func (f *Filesystem) LsRespond(args spec.LsArgs, success *bool) error {
	storeRWMutex.RLock()
	defer storeRWMutex.RUnlock()
	_, ok := store[args.Filename]
	if ok {
		*success = true
		return nil
	}

	return errors.New(fmt.Sprintf("File [%s] not found.", args.Filename))
}
