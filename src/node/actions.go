// Actual code implementations of RPCs
package node

import (
	"log"

	"../spec"
)

// Store the filename and data on this machine
// Also dispatch RPC calls to replica nodes
func _putAssign(args *spec.PutArgs) {
	log.SetPrefix(log.Prefix() + "_putAssign(): ")
	defer log.SetPrefix(spec.Prefix)
	log.Println(args)
}
