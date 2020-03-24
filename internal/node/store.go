package node

import (
	"fmt"
	"net/rpc"

	"github.com/slin63/chord-dfs/internal/config"
	"github.com/slin63/chord-dfs/internal/spec"
)

// Typical flow:
//   1. Client passes entry to Raft leader
//   2. Raft leader replicates entry and passes back to Filesystem.Execute RPC

// Store (function to initiate STORE action) (from: execute/client)
//   - Return the keys to the store for some given PID
func Store(args *spec.StoreArgs) ([]string, error) {
	var client *rpc.Client
	var err error
	PID := args.PID
	var storeKeys []string

	// Either:
	//   - Get store from own server
	//   - Get store from other server
	if PID == self.PID {
		return getStoreKeys(), nil
	}

	// Failed to connect
	client, err = connectTimeout(PID, config.C.RPCTimeout)
	if err != nil {
		config.LogIf(
			fmt.Sprintf("[STORE-X] [PID=%d] Not responding.", PID),
			config.C.LogStore,
		)
		return storeKeys, err
	}
	defer client.Close()

	// Connected successfully
	storeKeys, err = callStore(PID, args, client)
	if err != nil {
		config.LogIf(
			fmt.Sprintf("[STORE-X] [PID=%d] Error received from callStore().", PID),
			config.C.LogStore,
		)
		return storeKeys, err
	}

	return storeKeys, nil
}

// callStore (initiates StoreRespond RPC on given PID)
func callStore(PID int, args *spec.StoreArgs, client *rpc.Client) ([]string, error) {
	var storeKeys []string
	if err := (*client).Call("Filesystem.StoreRespond", *args, &storeKeys); err != nil {
		return storeKeys, err
	}
	return storeKeys, nil
}

// StoreRespond (receive RPC) (from: another server)
// Try and delete the given file
func (f *Filesystem) StoreRespond(args spec.StoreArgs, storeKeys *[]string) error {
	*storeKeys = getStoreKeys()
	return nil
}

func getStoreKeys() []string {
	storeRWMutex.RLock()
	defer storeRWMutex.RUnlock()

	storeKeys := make([]string, 0, len(store))
	for file := range store {
		storeKeys = append(storeKeys, file)
	}
	return storeKeys
}
