package calcstats

import "sync"

var Cnt int64
var lock = sync.Mutex{}

func Incr() {
	lock.Lock()
	Cnt++
	lock.Unlock()
}

var Tpose int64
var tlock = sync.Mutex{}

func TposeIncr() {
	tlock.Lock()
	Tpose++
	tlock.Unlock()
}