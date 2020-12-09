package model

type GlobalSocketInfo struct {
	Count         int
	DataChannel   chan []byte
	SignalChannel chan int
}
