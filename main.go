package main

import (
	cc "channelcontrol/endpoints"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/batrace/canceltoken"
)

func main() {
	ct := canceltoken.NewCancelToken()

	c := make(chan int, 100)
	sender := cc.NewSender(c)
	receiver := cc.NewReceiver(c)

	startReceiver(ct, receiver)
	startSender(ct, sender)

	sampleRates(ct, sender, receiver)

	ct.Wait()
}

func startReceiver(t *canceltoken.CancelToken, r *cc.Receiver) {
	t.Add(1)
	go func(t *canceltoken.CancelToken, r *cc.Receiver) {
		defer t.Done()
		for !t.IsCancelled() {
			r.Receive()
		}
		fmt.Println("Receiver exit")
	}(t, r)
}

func startSender(t *canceltoken.CancelToken, s *cc.Sender) {
	t.Add(1)
	go func(t *canceltoken.CancelToken, s *cc.Sender) {
		defer t.Done()
		defer s.Close()
		for !t.IsCancelled() {
			s.Send()
		}
		fmt.Println("Sender exit")
	}(t, s)
}

func sampleRates(t *canceltoken.CancelToken, s *cc.Sender, r *cc.Receiver) {
	var totalSent uint64
	var totalReceived uint64
	t.Add(1)
	go func(t *canceltoken.CancelToken, s *cc.Sender, r *cc.Receiver) {
		defer t.Done()
		for !t.IsCancelled() {
			duration := 1 * time.Second
			time.Sleep(duration)

			sc := atomic.SwapUint64(&s.SendCounter, 0)
			totalSent += sc
			rc := atomic.SwapUint64(&r.ReceiveCounter, 0)
			totalReceived += rc

			sendRate := float64(sc) / (duration.Seconds())
			receiveRate := float64(rc) / (duration.Seconds())

			fmt.Printf("SendRate = %f ReceiveRate = %f\n", sendRate, receiveRate)
		}
		fmt.Printf("Sent = %d Received = %d\n", totalSent, totalReceived)
	}(t, s, r)
}
