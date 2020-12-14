package log_watcher_pipe

import (
	"errors"
	"sync"
	"time"
)

type LogWatcherPipe interface {
	new(innerChannel <-chan []byte, outerChannels map[string]chan<- []byte) error
	Watching()
	Close()
	InsertOuterChannel(string, chan<- []byte) error
	IsInPipe(string) bool
}

type logWatcherPipe struct {
	inner        <-chan []byte
	outers       map[string]chan<- []byte
	transferLock *sync.Mutex
	isWatching   bool
}

func (p *logWatcherPipe) new(inner <-chan []byte, outers map[string]chan<- []byte) error {
	if inner == nil || outers == nil || len(outers) == 0 {
		return errors.New("inner channel or outer channels is nil")
	}
	p.inner = inner
	p.outers = outers
	p.transferLock = new(sync.Mutex)
	return nil
}

func (p logWatcherPipe) Watching() {
	if p.isWatching {
		return
	}
	go func() {
		p.isWatching = true
		for p.inner != nil && len(p.outers) > 0 {
			data := <-p.inner
			p.transferLock.Lock()
			t := make(chan int, 10)
			wg := new(sync.WaitGroup)
			for s, c := range p.outers {
				outerName := s
				outer := c
				wg.Add(1)
				go func() {
					t <- 1
					defer wg.Done()
					select {
					case outer <- data:
						break
					case <-time.After(time.Second * 5):
						delete(p.outers, outerName)
						break
					}
					<-t
				}()
			}
			wg.Wait()
			p.transferLock.Unlock()
		}
	}()
}

func (p *logWatcherPipe) Close() {
	p.transferLock.Lock()
	for s := range p.outers {
		delete(p.outers, s)
	}
	p.transferLock.Unlock()
}

func (p *logWatcherPipe) InsertOuterChannel(outerName string, outer chan<- []byte) error {
	p.transferLock.Lock()
	defer p.transferLock.Unlock()
	if _, isExist := p.outers[outerName]; isExist {
		return errors.New("outer channel is exist")
	}
	p.outers[outerName] = outer
	return nil
}

func (p logWatcherPipe) IsInPipe(outerName string) bool {
	_, ok := p.outers[outerName]
	return ok
}

func NewLogWatcherPipe(inner <-chan []byte, outers map[string]chan<- []byte) (LogWatcherPipe, error) {
	pipe := new(logWatcherPipe)
	err := pipe.new(inner, outers)
	if err != nil {
		return nil, err
	}
	return pipe, nil
}
