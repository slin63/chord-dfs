// Some custom receivers for a channel semaphore to facilitate logging.
package sem

import "log"

type empty struct{}
type Semaphore chan empty

const logging = false

func (s Semaphore) Lock() {
	if logging {
		log.Printf("Semaphore.Lock(): Locking!")
	}
	s <- empty{}
}

func (s Semaphore) Unlock() {
	if logging {
		log.Printf("Semaphore.Unlock(): Unlocking!")
	}
	if len(s) == 0 {
		panic("Semaphore.Unlock() called on unlocked semaphore")
	}
	<-s
}
