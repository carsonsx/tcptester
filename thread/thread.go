package thread

var signal = make(chan bool)

func Wait()  {
	<-signal
}

func Notify()  {
	signal <- false
}
