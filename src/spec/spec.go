package spec

import (
	"log"
	"net/rpc"
	"sort"
	"time"

	"../sem"
)

// Semaphores
var SelfSem = make(sem.Semaphore, 1)

// Logging prefix
const Prefix = "[DFS] - "

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

type Self struct {
	M            int
	PID          int
	MemberMap    MemberMapT
	FingerTable  FingerTableT
	SuspicionMap SuspicionMapT
}

// DFS RPCs
type PutArgs struct {
	Filename string
	Data     []byte
}

const FilesystemRPCPort = "6003"

const MemberRPCPort = "6002"
const MemberRPCRetryInterval = 2
const MemberRPCRetryMax = 5
const MemberInterval = 5

func ReportOnline(selfPID int) {
	log.Printf("[ONLINE] [PID=%d]", selfPID)
}

// Find the nearest PID to the given FPID on the virtual ring
// (including this node's own PID)
func GetSuccPID(FPID int, self *Self) *int {
	SelfSem.Lock()
	PIDs := []int{}
	for PID, _ := range (*self).MemberMap {
		PIDs = append(PIDs, PID)
	}
	SelfSem.Unlock()
	sort.Ints(PIDs)
	diff := 10000
	var succPID int
	FPID = FPID % (1 << self.M)

	// Find the smallest (FPID - PID) that is (> 0)
	// in an ordered array of ints
	for i := 0; i < len(PIDs); i++ {
		log.Println(PIDs[i])
		iterdiff := PIDs[i] - FPID
		if (iterdiff) < diff && iterdiff > 0 {
			diff = iterdiff
			succPID = PIDs[i]
		}
	}
	return &succPID
}

// Query the membership service running on the same machine for membership information.
func GetSelf() Self {
	var client *rpc.Client
	var err error
	for i := 0; i <= MemberRPCRetryMax; i++ {
		time.Sleep(MemberRPCRetryInterval * time.Second)
		client, err = rpc.DialHTTP("tcp", "localhost:"+MemberRPCPort)
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
	return reply
}
