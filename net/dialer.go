package net

import (
	"context"
	"github/go-robot/protocols"
	"sync"
)

type MyDialer interface {
	Disconnect()
	SendPacket(data []byte) bool
	ReadPacket() <-chan *protocols.Protocol
	Run(ctx context.Context, wg *sync.WaitGroup) bool
}

