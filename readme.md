# chord-dfs
A Chord based distributed file system.

## Implementation
TODO: Link blog post here.

## Setup
#### Launching with `docker-compose`
1. `docker-compose build && docker-compose up --remove-orphans --scale worker=<worker_count>`
    - Start 1 + `worker_count` nodes.
    - Recommended `worker_count ~= 5`. CPU utilization is high across all three components so expect some sluggishness.
2. Start client with `CONFIG=$(pwd)/config/config.dfs.json go run ./cmd/client/main.go <command>`. Available client commands listed below in _Client Commands_.

#### Configuration
`config.` files for each component can be found inside `/config`. Mappings are as follows:
- `config.dfs.json`: _Distributed File System Layer_
- `config.fd.json`: _Membership/Failure Detection Layer_
- `config.raft.json`: _Consensus Layer_

## Useage
#### Client Commands
1. `put localfilename sdfsfilename` (from local dir)
    - `put` both inserts _and_ updates a file
2. `get sdfsfilename localfilename` (fetches to local dir)
3. `delete sdfsfilename`
4. `ls filename` (list all machines where this data is stored)
5. `store` (list all files stored on this machine)
