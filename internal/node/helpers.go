package node

import (
    "log"
    "net"
    "net/http"
    "net/rpc"

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

// Connect to some RPC server and return a pointer to the client
func connect(PID int) *rpc.Client {
    node, ok := self.MemberMap[PID]
    if !ok {
        log.Fatalf("[CONNERR-X] Node with [PID=%d] not found. Exiting.", PID)
    }
    client, err := rpc.DialHTTP("tcp", (*node).IP+":"+config.C.FilesystemRPCPort)
    if err != nil {
        log.Fatal("put() dialing:", err)
    }

    return client
}
