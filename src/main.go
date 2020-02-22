package main

import (
	"log"

	"./node"
)

const logf = "dfs.log"
const prefix = "[DFS] - "

func main() {
	log.SetPrefix(prefix)
	// isIntroducer, err := strconv.ParseBool(os.Getenv("INTRODUCER"))
	// if err != nil {
	// 	log.Fatal("INTRODUCER not set in this environment")
	// }
	node.Live(logf)
}
