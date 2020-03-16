# chord-dfs
A Chord based distributed file system.

### Implementation
TODO: Link blog post here.

### Setup

1. `docker-compose build && docker-compose up --remove-orphans --scale worker=3` OR `docker build --tag foo . && docker run --rm -e INTRODUCER=1 foo`
2. Start client with `CONFIG=$(pwd)/config.json go run ./cmd/dfs-client/main.go <command>`. Available client commands listed below in _Client Commands_.

#### Client Commands
1. `put localfilename sdfsfilename` (from local dir)
    - `put` both inserts _and_ updates a file
2. `get sdfsfilename localfilename` (fetches to local dir)
3. `delete sdfsfilename`
4. `ls filename` (list all machines where this data is stored)
5. `store` (list all files stored on this machine)
