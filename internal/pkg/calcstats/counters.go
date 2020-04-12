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

var Dots int64
var dlock = sync.Mutex{}

func Dot() {
	dlock.Lock()
	Dots++
	dlock.Unlock()
}

var Crosses int64
var clock = sync.Mutex{}

func Cross() {
	clock.Lock()
	Crosses++
	clock.Unlock()
}

var Ns int64
var nlock = sync.Mutex{}

func Normalize() {
	nlock.Lock()
	Ns++
	nlock.Unlock()
}
