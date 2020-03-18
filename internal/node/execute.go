// RPC that executes arbitrary entries from the Raft Leader
package node

import (
    "fmt"
    "log"
    "strings"

    "github.com/slin63/chord-dfs/internal/config"
    "github.com/slin63/chord-dfs/internal/spec"
    "github.com/slin63/chord-dfs/pkg/parser"
    "github.com/slin63/raft-consensus/pkg/responses"
)

// (RPC from: Consensus Layer). Executes the given entry and returns a result.
func (f *Filesystem) Execute(entry string, result *responses.Result) error {
    config.LogIf(fmt.Sprintf("[EXECUTE] Executing %s", tr(entry, 20)), config.C.LogExecute)
    method, args, ok := parser.ParseEntry(strings.Split(entry, " "))
    // If this isn't a valid entry, explode. Something is really wrong upstream.
    if !ok {
        log.Fatalf("[EXECUTE-X] Invalid entry %s", tr(entry, 20))
    }

    // Actually execute the method with the arguments as identified by parser.ParseEntry
    execute(method, args, result)
    *result = responses.Result{Entry: entry, Success: true}
    return nil
}

func execute(method parser.MethodType, args []string, result *responses.Result) error {
    methodS, _ := parser.MethodString(method)
    var sdfs string = args[1]
    var b []byte = []byte(args[2])
    switch method {
    case parser.PUT:
        Put(&spec.PutArgs{sdfs, b, spec.NILPID})
    default:
        panic("TODO: Add the other methods")
    }
    config.LogIf(fmt.Sprintf("[EXECUTE] [METHOD=%s] [ARGS=%v %v %v]", methodS, args[0], args[1], tr(args[2], 20)), config.C.LogExecute)

    return nil
}
