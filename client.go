package main

import (
	"fmt"
	"github/go-robot/protocols"
)

type Client struct {
	serial uint32
	PlatformData

}

func NewClient(index uint32, pb PlatformData) *Client {
	c := new(Client)
	c.serial = index
	c.PID = pb.PID
	c.Nickname = pb.Nickname
	return c
}

func (c *Client)Update(ch chan<- []byte) {
	//fmt.Printf("-----------%d update\n", c.serial)
}

func (c *Client)OnConnected(ch chan<- []byte)  {
	ch<- protocols.PacketPing()
	fmt.Printf("-----------%d connected\n", c.serial)
}

func (c *Client)OnDisconnected()  {
	fmt.Printf("-----------%d disconnected\n", c.serial)
}

func (c *Client)ProcessProtocols(ch chan<- []byte, pbBuff []byte) {
	fmt.Println(pbBuff)
}
