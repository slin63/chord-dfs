package main

import (
	"log"
	"os"
	"strconv"

	"./client"
	"./node"
	"./spec"
)

const logf = "dfs.log"

func main() {
	log.SetPrefix(spec.Prefix + " - ")
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
