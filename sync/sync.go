package sync

import (
	"github.com/carsonsx/tcptester/util"
	"time"
	"github.com/carsonsx/log4g"
)

var _functions = util.NewQueue()

func Call(functions ...func()) {
	for _, function := range functions {
		_functions.Offer(function)
	}
	Done()
}

var done = make(chan bool)
var closed bool

func Wait(delay ...int) {
	<-done
	if len(delay) > 0 && delay[0] > 0 {
		time.Sleep(time.Duration(delay[0]) * time.Second)
	}
	closed = true
}

func Done(force ...bool) {
	if len(force) > 0 && force[0] {
		done <- true
	} else if _functions.Len() > 0 {
		_functions.Poll().(func())()
	} else {
		log4g.Info("all calls done")
		if !closed {
			done <- true
		}
	}
}
