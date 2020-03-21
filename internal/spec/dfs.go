package spec

type WriteCmd struct {
    Name string
    Data []byte
}

// DFS RPCs
type PutArgs struct {
    // File data
    Filename string
    Data     []byte

    // Server data
    From int

    // Whether or not this machine should replicate this file
    Replicate bool
}

type GetArgs struct {
    // File data
    Filename string
}
