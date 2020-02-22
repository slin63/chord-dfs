package spec

import (
	"log"
	"net/rpc"
	"time"
)

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
	PID          int
	MemberMap    MemberMapT
	FingerTable  FingerTableT
	SuspicionMap SuspicionMapT
}

const RPCPort = "6002"
const RPCRetryInterval = 2
const RPCRetryMax = 5
const MembershipInterval = 5

func ReportOnline(selfPID int) {
	log.Printf("[ONLINE] [PID=%d]", selfPID)
}

// Query the membership service running on the same machine for membership information.
func GetSelf() Self {
	var client *rpc.Client
	var err error
	for i := 0; i <= RPCRetryMax; i++ {
		time.Sleep(RPCRetryInterval * time.Second)
		client, err = rpc.DialHTTP("tcp", "localhost:"+RPCPort)
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
