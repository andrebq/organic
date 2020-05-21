package testutil

import (
	"runtime"
	"sync/atomic"
	"testing"

	"github.com/andrebq/organic/cell"
	"github.com/andrebq/organic/helper/signal"
)

// ExchangeMessage runs a simple message exchange between two membranes
// created by a medium
func ExchangeMessage(t *testing.T, prefix string, m cell.Medium) {
	c1 := cell.Grow("exchange_message_"+prefix+"_1", m)
	c2 := cell.Grow("exchange_message_"+prefix+"_2", m)

	sent := "abc123"
	var exit int32
	go func() {
		for atomic.LoadInt32(&exit) != 1 {
			c1.Emit(c2.ID(), "R", signal.Build().Str(sent).Signal())
			runtime.Gosched()
		}
	}()

	from, receptor, sig := c2.Recv()
	atomic.StoreInt32(&exit, 1)

	if receptor != "R" {
		t.Errorf("Invalid receptor should be [%v] got [%v]", "R", receptor)
	}

	if from != c1.ID() {
		t.Errorf("Should have received from [%v] got [%v]", c1.ID(), from)
	}

	var recv string
	signal.Match(sig, signal.String(&recv))
	if recv != "abc123" {
		t.Errorf("Invalid value, should be [%v] got [%v]", sent, recv)
	}
}
