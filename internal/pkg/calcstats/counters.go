package calcstats

import "sync"

var Cnt int64
var lock = sync.Mutex{}

func Incr() {
	lock.Lock()
	Cnt++
	lock.Unlock()
}
