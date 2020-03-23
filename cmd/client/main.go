// Client for interfacing with DFS through the Raft leader
package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "os/exec"
    "strings"

    "github.com/slin63/chord-dfs/internal/client"
    "github.com/slin63/chord-dfs/internal/config"
)

func main() {
    log.SetPrefix(config.C.Prefix + " - ")
    // Initialize logging to file
    f, err := os.OpenFile(config.C.LogfileClient, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("error opening file: %v", err)
    }
    defer f.Close()

    s := bufio.NewScanner(os.Stdin)
    w := os.Stdout
    fmt.Fprint(w, client.HelpS+"\n")
    for {
        fmt.Fprint(w, "> ")
        s.Scan() // get next the token
        switch strings.ToLower(s.Text()) {
        case "!ls":
            b, err := exec.Command("ls").Output()
            if err != nil {
                log.Fatal(err)
            }
            fmt.Fprint(w, string(b))
        case "":
            continue
        default:
            client.Parse(strings.Split(s.Text(), " "))
        }

        // todo:
        // don't crash on errors
        // handle arrow key inputs and normal editing
    }

    client.Parse(os.Args[1:])
}
