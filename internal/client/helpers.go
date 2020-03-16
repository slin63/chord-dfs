package client

import (
    "log"
    "strings"
)

const helpS = `Available operations:
1. put localfilename sdfsfilename (from local dir)
2. get sdfsfilename localfilename (fetches to local dir)
3. delete sdfsfilename
4. ls filename (list all machines where this data is stored)
5. store (list all files stored on this machine)`

// Make sure that the entry is a valid command
func parseEntry(args []string) (string, bool) {
    switch args[0] {
    case "put":
        if len(args) != 3 {
            return "", false
        }
        return strings.Join(args, " "), true
    case "get":
        log.Println("get")
        return strings.Join(args, " "), true
    case "delete":
        log.Println("delete")
        return strings.Join(args, " "), true
    default:
        return "", false
    }
}
