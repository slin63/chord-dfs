package node

import (
    "errors"
    "log"
    "net"
    "net/http"
    "net/rpc"
    "time"

    "github.com/slin63/chord-dfs/internal/config"
)

func serveFilesystemRPC() {
    fs := new(Filesystem)
    rpc.Register(fs)
    rpc.HandleHTTP()
    l, e := net.Listen("tcp", ":"+config.C.FilesystemRPCPort)
    if e != nil {
        log.Fatal("[ERROR] serveFilesystemRPC():", e)
    }
    http.Serve(l, nil)
}

func connectTimeout(PID, timeout int) (*rpc.Client, error) {
    var client *rpc.Client
    c := make(chan *rpc.Client)
    e := make(chan error)
    go func() {
        client, err := connect(PID)
        if err != nil {
            e <- err
            return
        }
        c <- client
    }()
    select {
    case client := <-c:
        return client, nil
    case err := <-e:
        return client, err
    case <-time.After(time.Duration(timeout) * time.Second):
        return client, errors.New("Timed out waiting for response")
    }
}

// Connect to some RPC server and return a pointer to the client
func connect(PID int) (*rpc.Client, error) {
    node, ok := self.MemberMap[PID]
    if !ok {
        log.Fatalf("[CONNERR-X] Node with [PID=%d] not found. Exiting.", PID)
    }
    client, err := rpc.DialHTTP("tcp", (*node).IP+":"+config.C.FilesystemRPCPort)
    if err != nil {
        log.Println("put() dialing:", err)
        return client, err
    }

    return client, nil
}

func tr(s string, max int) string {
    if len(s) > max {
        return s[0:max] + "..."
    }
    return s
}

func trb(b []byte, max int) string {
    return tr(string(b), max)
}
