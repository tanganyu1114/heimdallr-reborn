package log_watcher_pipe

import (
	"fmt"
	"testing"
	"time"
)

func TestNewLogWatcherBuffer(t *testing.T) {
	var str string
	w := make(chan []byte, 1)
	r1 := make(chan []byte, 1)
	r2 := make(chan []byte, 1)
	r3 := make(chan []byte, 1)
	pipe, err := NewLogWatcherPipe(w, map[string]chan<- []byte{"test1": r1})
	if err != nil {
		t.Fatal(err)
		return
	}
	pipe.Watching()
	go func() {
		for {
			s := <-r1
			fmt.Println("test1:", string(s))
		}
	}()
	go func() {
		for {
			s := <-r2
			fmt.Println("test2:", string(s))
		}
	}()
	go func() {
		for {
			s := <-r3
			fmt.Println("test3:", string(s))
			time.Sleep(time.Second * 10)
		}
	}()
	w <- []byte("test111")
	w <- []byte("test222")
	time.Sleep(time.Second * 2)
	err = pipe.InsertOuterChannel("test2", r2)
	if err != nil {
		t.Fatal(err)
		return
	}
	str = "111111"
	w <- []byte(str)
	time.Sleep(time.Second * 5)
	w <- []byte("333333")
	err = pipe.InsertOuterChannel("test3", r3)
	w <- []byte("4444444")
	w <- []byte("end")
	time.Sleep(time.Second * 4)
	pipe.Close()
}
