// RPC that executes arbitrary entries from the Raft Leader
package node

import (
    "fmt"

    "github.com/slin63/chord-dfs/internal/config"
    "github.com/slin63/raft-consensus/pkg/responses"
)

func (f *Filesystem) Execute(entry string, result *responses.Result) error {
    config.LogIf(fmt.Sprintf("[EXECUTE] Executing %s", entry), config.C.LogExecute)
    config.LogIf(fmt.Sprintf("[EXECUTE] TODO", entry), config.C.LogExecute)
    *result = responses.Result{Entry: entry, Success: true}
    return nil
}
