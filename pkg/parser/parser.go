// DFS-specific parsing of arguments
package parser

import (
    "log"
    "strconv"
    "strings"
)

type MethodType int

const (
    PUT MethodType = iota
    GET
    DELETE
)

func MethodString(method MethodType) (string, bool) {
    switch method {
    case PUT:
        return "PUT", true
    case GET:
        return "GET", true
    case DELETE:
        return "DELETE", true
    default:
        return "", false
    }
}

// Make sure that the entry is a valid command
// Available operations:
//   1. put localfilename sdfsfilename (from local dir)
//   2. get sdfsfilename localfilename (fetches to local dir)
//   3. delete sdfsfilename
//   4. ls filename (list all machines where this data is stored)
//   5. store (list all files stored on this machine)
func ParseEntry(args []string) (MethodType, []string, bool) {
    cleanTerm(args)
    method := args[0]
    args = args[1:]
    switch strings.ToLower(method) {
    case "put":
        if len(args) < 2 {
            return 0, args, false
        }
        return PUT, cleanData(args, 2), true
    case "get":
        if len(args) < 2 {
            return 0, args, false
        }
        return GET, args, true
    case "delete":
        log.Println("delete")
        return 0, args, false
    default:
        return 0, args, false
    }
}

// Sometimes file data that has whitespace in it is split into multiple arguments.
// We can use this function to rejoin that split data.
//   start :: the first index of the data that was wrongly split into multiple indices
func cleanData(args []string, start int) []string {
    if len(args) > start {
        args[start] = strings.Join(args[start:], " ")
        args = args[0 : start+1]
    }
    return args
}

// Return args without a term if a term is present
func cleanTerm(args []string) {
    first := args[0]
    term := string(first[0])
    var hasTerm bool
    if _, err := strconv.Atoi(term); err == nil {
        hasTerm = true
    }

    if hasTerm {
        args[0] = first[2:]
    }
}
