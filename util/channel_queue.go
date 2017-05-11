package util

type chqueue struct {
	size   int
	data   chan interface{}
	length int
}

func NewChannelQueue(size int) *chqueue {
	var q chqueue
	q.size = size
	q.data = make(chan interface{}, size)
	return &q
}

func (q *chqueue) Offers(a []interface{}) {
	for _, e := range a {
		q.Offer(e)
	}
}

func (q *chqueue) Offer(v interface{}) {
	q.data <- v
	q.length++
}

func (q *chqueue) Poll() interface{} {
	if q.length == 0 {
		return nil
	}
	q.length--
	return <-q.data
}

func (q *chqueue) Clear() {
	for q.Len() > 0 {
		<-q.data
	}
	q.length = 0
}

func (q *chqueue) Size() int {
	return q.size
}

func (q *chqueue) Len() int {
	return q.length
}
