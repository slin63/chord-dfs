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
func ParseEntry(args []string) (MethodType, []string, bool) {
    cleanTerm(args)
    switch strings.ToLower(args[0]) {
    case "put":
        if len(args) < 3 {
            return 0, args[1:], false
        }
        return PUT, args[1:], true
    case "get":
        log.Println("get")
        return 0, args[1:], false
    case "delete":
        log.Println("delete")
        return 0, args[1:], false
    default:
        return 0, args[1:], false
    }
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
