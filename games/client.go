package games

import "github/go-robot/protocols"

type Client interface {
	Update(ch chan<- []byte)
	OnConnected(ch chan<- []byte)
	OnDisconnected()
	ProcessProtocols(ch chan<- []byte, p *protocols.Protocol) bool
}

