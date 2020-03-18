// RPC that executes arbitrary entries from the Raft Leader
package node

import (
    "fmt"
    "strings"

    "github.com/slin63/chord-dfs/internal/config"
    "github.com/slin63/chord-dfs/pkg/parser"
    "github.com/slin63/raft-consensus/pkg/responses"
)

// TODO (03/17 @ 17:37): Parse the entry and deal with it appropriately
func (f *Filesystem) Execute(entry string, result *responses.Result) error {
    config.LogIf(fmt.Sprintf("[EXECUTE] Executing %s", entry), config.C.LogExecute)
    method, args, ok := parser.ParseEntry(strings.Split(entry, " "))
    methodS, _ := parser.MethodString(method)
    if ok {
        // Do some stuff
        config.LogIf(fmt.Sprintf("[EXECUTE] [METHOD=%s] [ARGS=%v]", methodS, args), config.C.LogExecute)
    }
    *result = responses.Result{Entry: entry, Success: true}
    return nil
}
