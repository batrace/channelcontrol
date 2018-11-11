package channelcontrol

import "sync/atomic"

// Receiver receives on a channel
type Receiver struct {
	ReceiveCounter uint64
	c              chan int
}

// NewReceiver creates a new receiver
func NewReceiver(channel chan int) *Receiver {
	r := Receiver{}
	r.c = channel
	return &r
}

// Receive receives data until the channel is closed
func (r *Receiver) Receive() {
	for range r.c {
		atomic.AddUint64(&r.ReceiveCounter, 1)
	}
}
