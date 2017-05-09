package sync

import "github.com/carsonsx/tcptester/util"

var _functions = util.NewQueue()

func Call(functions ...func()) {
	for _, function := range functions {
		_functions.Offer(function)
	}
	Done()
}

var done = make(chan bool)

func Wait() {
	<-done
}

func Done() {
	if _functions.Len() > 0 {
		_functions.Poll().(func())()
	} else {
		done <- true
	}
}

func ForceDone() {
	done <- true
}
