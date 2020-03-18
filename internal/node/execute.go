// RPC that executes arbitrary entries from the Raft Leader
package node

import (
    "fmt"
    "io/ioutil"
    "log"
    "strings"

    "github.com/slin63/chord-dfs/internal/config"
    "github.com/slin63/chord-dfs/internal/spec"
    "github.com/slin63/chord-dfs/pkg/parser"
    "github.com/slin63/raft-consensus/pkg/responses"
)

// (RPC from: Consensus Layer). Executes the given entry and returns a result.
func (f *Filesystem) Execute(entry string, result *responses.Result) error {
    config.LogIf(fmt.Sprintf("[EXECUTE] Executing %s", entry), config.C.LogExecute)
    method, args, ok := parser.ParseEntry(strings.Split(entry, " "))
    // If this isn't a valid entry, explode. Something is really wrong upstream.
    if !ok {
        log.Fatalf("[EXECUTE-X] Invalid entry %s", entry)
    }

    // Actually execute the method with the arguments as identified by parser.ParseEntry
    execute(method, args, result)
    *result = responses.Result{Entry: entry, Success: true}
    return nil
}

func execute(method parser.MethodType, args []string, result *responses.Result) error {
    methodS, _ := parser.MethodString(method)
    var local string = args[0]
    var sdfs string = args[1]
    switch method {
    case parser.PUT:
        Put(putArgs(local, sdfs))
    default:
        panic("TODO: Add the other methods")
    }
    config.LogIf(fmt.Sprintf("[EXECUTE] [METHOD=%s] [ARGS=%v]", methodS, args), config.C.LogExecute)

    return nil
}

//
// - Read file `local`
// - Send RPC to DFS server with correct arguments
// - DFS server decides what to do with it and where to put it
func putArgs(local, sdfs string) *spec.PutArgs {
    f, err := ioutil.ReadFile(local)
    if err != nil {
        log.Fatal("[putArgs()]: ", err)
    }
    return &spec.PutArgs{sdfs, f, spec.NILPID}
}
