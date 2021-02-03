package main

type Client struct {
	serial int
	pid int32


}

func NewClient(index int) *Client {
	c := new(Client)
	c.serial = index
	return c
}

func (c *Client)Update() {

}

func (c *Client)ProcessProtocols(ch chan<- []byte, pbBuff []byte) {

}
