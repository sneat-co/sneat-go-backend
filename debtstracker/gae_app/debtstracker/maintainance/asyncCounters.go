package maintainance

//import (
//	"sync"
//
//	"github.com/captaincodeman/datastore-mapper"
//)
//
//type asyncCounters struct {
//	sync.Mutex
//	locked   bool
//	counters mapper.Counters
//}
//
//func NewAsynCounters(counters mapper.Counters) *asyncCounters {
//	return &asyncCounters{counters: counters}
//}
//
//func (ac *asyncCounters) Increment(name string, delta int64) {
//	wasLocked := ac.locked
//	if !wasLocked {
//		ac.Lock()
//	}
//	ac.counters.Increment(name, delta)
//	if !wasLocked {
//		ac.Unlock()
//	}
//}
//
//func (ac *asyncCounters) Lock() {
//	ac.Mutex.Lock()
//	ac.locked = true
//}
//
//func (ac *asyncCounters) Unlock() {
//	if ac.locked {
//		ac.Mutex.Unlock()
//		ac.locked = false
//	}
//}
