// RPC that executes arbitrary entries from the Raft Leader
package node

import (
    "encoding/json"
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
    result.Entry = entry

    return nil
}

func execute(method parser.MethodType, args []string, result *responses.Result) error {
    methodS, _ := parser.MethodString(method)
    switch method {
    case parser.PUT:
        sdfs := args[1]
        b := []byte(args[2])
        replicas := Put(&spec.PutArgs{
            Filename:  sdfs,
            Data:      b,
            From:      spec.NILPID,
            Replicate: true,
        })
        data, err := json.Marshal(replicas)
        if err != nil {
            log.Fatal("[execute()] Error while marshaling response data:", err)
        }
        *result = responses.Result{Success: true, Data: string(data)}
    case parser.GET:
        sdfs := args[0]
        data, err := Get(&spec.GetArgs{Filename: sdfs})
        *result = responses.Result{Success: true, Data: string(data)}
        if err != nil {
            *result = responses.Result{Success: false, Error: responses.FILENOTFOUND}
        }
    case parser.DELETE:
        sdfs := args[0]
        err := Delete(&spec.DeleteArgs{Filename: sdfs})
        *result = responses.Result{Success: true, Data: sdfs}
        if err != nil {
            *result = responses.Result{Success: false, Error: responses.FILENOTFOUND}
        }
    case parser.LS:
        sdfs := args[0]
        holders, err := Ls(&spec.LsArgs{Filename: sdfs})
        if err != nil {
            *result = responses.Result{Success: false, Error: responses.FILENOTFOUND}
        }

        data, err := json.Marshal(holders)
        if err != nil {
            log.Fatal("[execute()] Error while marshaling response data:", err)
        }
        *result = responses.Result{Success: true, Data: string(data)}
    default:
        panic("TODO: Add the other methods")
    }
    config.LogIf(fmt.Sprintf("[EXECUTE] [METHOD=%s] [ARGS=%s]", methodS, tr(strings.Join(args, " "), 40)), config.C.LogExecute)
    return nil
}
