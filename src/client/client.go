package client

import (
	"io/ioutil"
	"log"
)

// 9 MP3: Build SDFS
// 12   - Required operations:
// 1. `put localfilename sdfsfilename` (from local dir)
//   - `put` both inserts _and_ updates a file
// 2. `get sdfsfilename localfilename` (fetches to local dir)
// 3. `delete sdfsfilename`
// 4. `ls filename` (list all machines where this data is stored)
// 5. `store` (list all files stored on this machine)
const helpS = `Available operations:
1. put localfilename sdfsfilename (from local dir)
2. get sdfsfilename localfilename (fetches to local dir)
3. delete sdfsfilename
4. ls filename (list all machines where this data is stored)
5. store (list all files stored on this machine)`

const server = "localhost:6002"

func Parse(args []string) {
	if len(args) == 0 {
		log.Println(helpS)
		return
	}
	switch args[0] {
	case "put":
		if len(args) == 3 {
			put(args[1], args[2])
		}
	case "get":
		log.Println("get")
	case "delete":
		log.Println("delete")
	default:
		log.Println(helpS)
	}
}

// How to copy a file onto another server?
func put(local, sdfs string) {
	// - load in with read file
	// - send over tcp (have to use tcp because we're using rpcs) to server
	// - server decides what to do with it and where to put it
	f, err := ioutil.ReadFile(local)
	if err != nil {
		log.Fatal("put(): ", err)
	}
	log.Println(f)
}
