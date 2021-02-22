package net

import (
	"context"
	"github/go-robot/protocols"
	"log"
	"sync"
)

type MyDialer interface {
	Disconnect()
	SendPacket(data []byte) bool
	ReadPacket() <-chan *protocols.Protocol
	Run(ctx context.Context, wg *sync.WaitGroup) bool
}

const (
	WS = "ws"
)

func NewConnect(protocol, sAddr string) MyDialer {
	switch protocol {
	case WS:
		return NewWSConnect(sAddr)
	}
	log.Panic("game id is not undefined")
	return nil
}