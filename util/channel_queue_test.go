package util

import (
	"testing"
)

func TestChannelQueue(t *testing.T) {


	q := NewChannelQueue(10)
	q.Offer("a")
	q.Offer("b")
	q.Offer("c")


	if q.Poll() != "a" {
		t.FailNow()
	}
	q.Offer("d")

	if q.Poll() != "b" {
		t.FailNow()
	}

	if q.Poll() != "c" {
		t.FailNow()
	}

	if q.Poll() != "d" {
		t.FailNow()
	}

	if q.Len() != 0 {
		t.Error(q.Len())
		t.FailNow()
	}

}
