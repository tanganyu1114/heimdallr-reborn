package log_watcher_buffer

import (
	"testing"
	"time"
)

func TestNewLogWatcherBuffer(t *testing.T) {
	lb := NewLogWatcherBuffer()
	r1 := make(chan []byte, 1)
	r2 := make(chan []byte, 1)
	w := make(chan []byte, 1)
	go func() {
		lb.Read(r1)
	}()
	go func() {
		lb.Read(r2)
	}()
	go func() {
		for {
			lb.Write(w)
		}
	}()
	go func() {
		for {
			select {
			case <-time.Tick(time.Second):
				lb.BufferCP()
			}
		}
	}()
	w <- []byte("test111")
	w <- []byte("test222")
	t.Log(string(<-r1))
	t.Log(string(<-r2))
}
