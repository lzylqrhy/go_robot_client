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
	TCP = "tcp"
)

func NewConnect(protocol, sAddr string) MyDialer {
	switch protocol {
	case WS:
		return NewWSConnect(sAddr)
	case TCP:
		return NewTCPConnect(sAddr)
	}
	log.Panic("net protocol is not undefined")
	return nil
}