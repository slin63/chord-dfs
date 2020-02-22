package main

import (
	"log"
	"os"
	"strconv"

	"./client"
	"./node"
)

const logf = "dfs.log"
const prefix = "[DFS] - "

func main() {
	log.SetPrefix(prefix)
	server, err := strconv.ParseBool(os.Getenv("SERVER"))
	if err != nil {
		log.Fatal("SERVER not set in this environment")
	}
	if server {
		node.Live(logf)
	} else {
		client.Parse(os.Args[1:])
	}
}
