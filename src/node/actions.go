// Actual code implementations of RPCs
package node

import (
	"fmt"
	"log"

	"../spec"
)

// Store the filename and data on this machine
// Also dispatch RPC calls to replica nodes
func _putAssign(args *spec.PutArgs) {
	log.SetPrefix(log.Prefix() + "_putAssign(): ")
	defer log.SetPrefix(spec.Prefix + fmt.Sprintf(" [PID=%d]", self.PID) + " - ")
	log.Println(args)
}
