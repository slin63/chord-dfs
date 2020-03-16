package main

import (
    "log"
    "os"

    "github.com/slin63/chord-dfs/internal/client"
    "github.com/slin63/chord-dfs/internal/config"
)

func main() {
    log.SetPrefix(config.C.Prefix + " - ")
    // Initialize logging to file
    f, err := os.OpenFile(config.C.Logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("error opening file: %v", err)
    }
    defer f.Close()
    client.Parse(os.Args[1:])
}
