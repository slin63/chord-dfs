# ![](./images/fishsmall.png) Chord-ish DeFiSh
A Chord-based distributed file system that maintains state through a replicated log. Built using a bunch of other junk that I built. Explanation with graphics at the bottom!

Uses:
- [Chord-ish](https://github.com/slin63/chord-failure-detector#-chord-ish) as a membership and failure detection layer.
- [Leeky Raft](https://github.com/slin63/raft-consensus#-leeky-raft) as a consensus layer.

![](./images/startup.gif)

*Starting the Raft cluster.*

![](./images/leaderelection.gif)

*Electing a new leader on leader failure.*


## Setup
#### Launching with `docker-compose`
0. Setup the network.
    - `docker network create dfs-net`
1. `docker-compose build && docker-compose up --remove-orphans --scale worker=<worker_count>`
    - Start 1 + `worker_count` nodes.
    - Recommended `worker_count ~= 5`. CPU utilization is high across all three components so expect some sluggishness.
2. Build & run client with
    - `docker build --tag client . -f ./dockerfiles/client/Dockerfile; docker run --rm -it --network="dfs-net" client /bin/sh -c ./dfs`
    - `> put go.mod remote`.
    - `> get remote local`.
    -  Available client commands listed below in _Client Commands_.

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

## In a Nutshell

Chord-ish DeFiSh works by combining three separate layers, each of which I built from scratch and are coated in an alarming amount of my own blood, creaking from the rust that accumulated as a result of my tears and sweat getting all over them. They are listed in order of their role in the placing of a user's file onto the distributed filesystem.

1. [Chord-ish](https://github.com/slin63/chord-failure-detector#-chord-ish), the membership layer. The membership layer lays the foundation for everything by assigning nodes / servers in a network onto some "virtual ring", giving them a distinct ID number as a function of their IP address. Then each node begins heartbeating to some number of nodes around it, setting up a system that allows them to gossip membership information and become aware of failures.

2. [Leeky Raft](https://github.com/slin63/raft-consensus#-leeky-raft), the consensus layer. A client sends commands, or entries to the consensus layer. These commands are similar to HTTP verbs. For example, the command to put the file `test.txt` onto our distributed filesystem with the name `remote.txt` would be expressed as `"PUT test.txt remote.txt"`. The consensus layer then replicates this entry to all other nodes in the network. On confirmation that the replication was (mostly) successful, they send the command to the filesystem layer.

3. [Chordish DeFiSh](https://github.com/slin63/chord-dfs#-chord-ish-defish), the filesystem layer. The filesystem layer receives the command from the consensus layer and begins executing it. It assigns the file a distinct ID number as a function of their filename, using the same method as the membership layer. It then stores this file at the first node with an ID greater than or equal to its own. If no node's ID is greater, then it wraps around the ring and tries to find a node there.

   Files are replicated to the 2 nodes directly "ahead" of the aforementioned node. Files are stored as actual files in each nodes' filesystem, and as `filename:sha1(file data)` maps in the runtime memory of each Chordish DeFiSh process, as a fast way to check for file ownership & save time by ignoring write requests for a file it already has.

   From there, users can upload, delete, or download files from the file system. The visuals below will explain how this all works, sort of.

![](./images/1.jpg)

![](./images/2.jpg)

![](./images/3.jpg)

