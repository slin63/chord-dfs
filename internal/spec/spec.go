// Constants for configuration and dealing with the membership layer
package spec

import (
	"log"
	"net/rpc"
	"sort"
	"sync"
	"time"

	"github.com/slin63/chord-dfs/internal/config"
)

// Semaphores
var SelfRWMutex sync.RWMutex

// Membership RPCs
type MemberNode struct {
	// Address info formatted ip_address
	IP        string
	Timestamp int64
	Alive     bool
}
type Membership int
type SuspicionMapT map[int]int64
type FingerTableT map[int]int
type MemberMapT map[int]*MemberNode

const NILPID = -1

type Self struct {
	M            int
	PID          int
	MemberMap    MemberMapT
	FingerTable  FingerTableT
	SuspicionMap SuspicionMapT
}

// DFS RPCs
type PutArgs struct {
	// File data
	Filename string
	Data     []byte

	// Server data
	From int
}

func ReportOnline() {
	log.Printf("[ONLINE]")
}

// Find the nearest PID to the given FPID on the virtual ring
// (including this node's own PID)
func NearestPID(FPID int, self *Self) int {
	SelfRWMutex.RLock()
	PIDs := []int{}
	PIDsExtended := []int{}

	for PID := range (*self).MemberMap {
		PIDs = append(PIDs, PID)
	}
	SelfRWMutex.RUnlock()

	for _, PID := range PIDs {
		PIDsExtended = append(PIDsExtended, PID+(1<<self.M))
	}
	PIDs = append(PIDs, PIDsExtended...)

	sort.Ints(PIDs)
	diff := 10000
	var nearestPID int

	// Find the smallest (FPID - PID) that is (> 0)
	// in an ordered array of ints
	for i := 0; i < len(PIDs); i++ {
		iterdiff := PIDs[i] - FPID
		if (iterdiff < diff) && (iterdiff > 0) {
			diff = iterdiff
			nearestPID = PIDs[i] % (1 << self.M)
		}
	}
	return nearestPID
}

// Query the membership service running on the same machine for membership information.
func GetSelf(self *Self) {
	var client *rpc.Client
	var err error
	for i := 0; i <= config.C.MemberRPCRetryMax; i++ {
		time.Sleep(time.Duration(config.C.MemberRPCRetryInterval) * time.Second)
		client, err = rpc.DialHTTP("tcp", "localhost:"+config.C.MemberRPCPort)
		if err != nil {
			log.Println("RPC server still spooling... dialing:", err)
		} else {
			break
		}
	}
	// Synchronous call
	var reply Self
	err = client.Call("Membership.Self", 0, &reply)
	if err != nil {
		log.Fatal("RPC error:", err)
	}
	SelfRWMutex.Lock()
	*self = reply
	SelfRWMutex.Unlock()
}
