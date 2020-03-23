// Helper functions for file system operations
package filesys

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"

    "github.com/slin63/chord-dfs/internal/config"
)

// Destructive file writer
func Write(filename string, data []byte) int {
    err := ioutil.WriteFile(config.C.Filedir+filename, data, 0644)
    if err != nil {
        log.Fatal(err)
    }
    config.LogIf(
        fmt.Sprintf(
            "[WRITE] Wrote %s (%d bytes) to %s",
            filename,
            len(data),
            config.C.Filedir+filename),
        config.C.LogWrites,
    )
    return len(data)
}

func Read(filename string) []byte {
    data, err := ioutil.ReadFile(config.C.Filedir + filename)
    if err != nil {
        log.Fatal(err)
    }
    config.LogIf(
        fmt.Sprintf(
            "[READ] Read %s (%d bytes)",
            filename,
            len(data)),
        config.C.LogReads,
    )
    return data
}

func Remove(filename string) error {
    config.LogIf(
        fmt.Sprintf(
            "[DELETE] Removed %s",
            filename),
        config.C.LogDeletes,
    )
    return os.Remove(config.C.Filedir + filename)
}
