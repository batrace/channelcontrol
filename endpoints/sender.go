package channelcontrol

import "sync/atomic"

// Sender sends on a channel
type Sender struct {
	SendCounter uint64
	c           chan int
}

// NewSender creates a new sender
func NewSender(channel chan int) *Sender {
	s := Sender{}
	s.c = channel
	return &s
}

// Send sends data
func (s *Sender) Send() {
	s.c <- 1
	atomic.AddUint64(&s.SendCounter, 1)
}

// Close closes the sender channel
func (s *Sender) Close() {
	close(s.c)
}
