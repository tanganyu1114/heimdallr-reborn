package log_watcher_buffer

import (
	"bytes"
	"sync"
	"time"
)

type LogWatcherBuffer interface {
	Write(<-chan []byte)
	BufferCP()
	Read(chan<- []byte)
}

type logWatcherBuffer struct {
	writeBuf      *bytes.Buffer
	writeLock     *sync.Mutex
	currentLog    []byte
	readWaitGroup *sync.WaitGroup
	readLock      *sync.Mutex
	canRead       bool
}

func (lb *logWatcherBuffer) Write(bc <-chan []byte) {
	b := <-bc
	lb.writeLock.Lock()
	lb.writeBuf.Write(b)
	lb.writeLock.Unlock()
}

func (lb *logWatcherBuffer) BufferCP() {
	lb.writeLock.Lock()
	defer lb.writeLock.Unlock()
	lb.canRead = false
	c := make(chan bool, 1)
	go func(ch chan bool) {
		lb.readWaitGroup.Wait()
		c <- true
	}(c)
	select {
	case <-c:
		break
	case <-time.After(time.Second * 30):
		return
	}
	lb.readLock.Lock()
	lb.currentLog = lb.writeBuf.Bytes()
	lb.readLock.Unlock()
	lb.writeBuf.Reset()
	lb.canRead = true
}

func (lb *logWatcherBuffer) Read(bc chan<- []byte) {
	for {
		if !lb.canRead {
			time.Sleep(time.Millisecond)
			continue
		}
		lb.readLock.Lock()
		lb.readWaitGroup.Add(1)
		lb.readLock.Unlock()
		bc <- lb.currentLog
		lb.readWaitGroup.Done()
		return
	}
}

func NewLogWatcherBuffer() LogWatcherBuffer {
	return &logWatcherBuffer{
		writeBuf:      bytes.NewBuffer(make([]byte, 0, 1024)),
		writeLock:     new(sync.Mutex),
		currentLog:    make([]byte, 0, 1024),
		readWaitGroup: new(sync.WaitGroup),
		readLock:      new(sync.Mutex),
		canRead:       false,
	}
}
