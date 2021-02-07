package fish

import (
	"fmt"
	"github/go-robot/common"
	"github/go-robot/protocols"
	"strings"
)

type Client struct {
	index  uint32
	ptData *common.PlatformData
	sevTime uint32
}

func NewClient(index uint32, pd *common.PlatformData) *Client {
	c := new(Client)
	c.index = index
	c.ptData = pd
	fmt.Printf("new client: %v \n", c)
	fmt.Printf("new client: %v \n", *c.ptData)
	return c
}

func (c *Client)Update(ch chan<- []byte) {
	//fmt.Printf("client serial=%d update\n", c.serial)
}

func (c *Client)OnConnected(ch chan<- []byte)  {
	var ping protocols.C2SSyncTime
	ch<- ping.Bytes()
	fmt.Printf("client index=%d connected\n", c.index)
}

func (c *Client)OnDisconnected()  {
	fmt.Printf("client index=%d disconnected\n", c.index)
}

func (c *Client)ProcessProtocols(ch chan<- []byte, p *protocols.Protocol) bool {
	fmt.Printf("cmd:0x%04x, data: %v\n", p.Head.Cmd, p.Content)
	switch p.Head.Cmd {
	case protocols.SyncTimeCode:
		return c.ProcessSyncTime(ch, p)
	case protocols.S2CLoginCode:
		return c.ProcessLogin(ch, p)
	case protocols.EnterRoomCode:
		return c.ProcessEnterRoom(ch, p)
	case protocols.SceneInfoCode:
		return c.ProcessSceneInfo(ch, p)
	}
	return true
}

func (c *Client) ProcessSyncTime(ch chan<- []byte, p *protocols.Protocol) bool {
	s2cSync := new(protocols.S2CSyncTime)
	s2cSync.Parse(p)
	fmt.Println(s2cSync)
	c.sevTime = s2cSync.TimeStamp
	// 请求登录
	var s2cLogin protocols.C2SLogin
	s2cLogin.IsChildGame = false
	var strBuilder strings.Builder
	strBuilder.WriteString(c.ptData.LoginToken)
	strBuilder.WriteString(":0x20:1")
	s2cLogin.Token = strBuilder.String()
	fmt.Println("session:", s2cLogin.Token)
	ch<- s2cLogin.Bytes()
	return true
}

func (c *Client) ProcessLogin(ch chan<- []byte, p *protocols.Protocol) bool {
	var s2cLogin protocols.S2CLogin
	s2cLogin.Parse(p)
	if s2cLogin.Status == 1 {
		// 登录成功
		fmt.Printf("client index=%d, pid=%d login successfully\n", c.index, c.ptData.PID)
		// 发送资源加载完成
		var c2sLoaded protocols.C2SResourceLoaded
		ch<- c2sLoaded.Bytes()
		// 进入房间
		var c2sEnterRoom protocols.C2SEnterRoom
		c2sEnterRoom.RoomID = 20001
		ch<- c2sEnterRoom.Bytes()
		return true
	}
	fmt.Printf("client index=%d, pid=%d login failed, status: %d\n", c.index, c.ptData.PID, s2cLogin.Status)
	return false
}

func (c *Client) ProcessEnterRoom(ch chan<- []byte, p *protocols.Protocol) bool {
	var s2cEnterRoom protocols.S2CEnterRoom
	s2cEnterRoom.Parse(p)
	if s2cEnterRoom.Result != 0 {
		// 进入失败
		fmt.Printf("client index=%d, pid=%d enter room=%d failed, result: %d\n",
			c.index, c.ptData.PID, s2cEnterRoom.RoomID, s2cEnterRoom.Result)
		return false
	}
	fmt.Printf("client index=%d, pid=%d enter room=%d successfully\n", c.index, c.ptData.PID, s2cEnterRoom.RoomID)
	// 请求场景信息
	var c2sSceneInfo protocols.C2SGetSceneInfo
	ch<- c2sSceneInfo.Bytes()
	// 转发
	var c2sTransmit protocols.C2STransmitActivity
	c2sTransmit.Activity = "hello"
	ch<- c2sTransmit.Bytes()
	return true
}

func (c *Client) ProcessSceneInfo(ch chan<- []byte, p *protocols.Protocol) bool {
	var s2cSceneInfo protocols.S2CGetSceneInfo
	s2cSceneInfo.Parse(p)
	fmt.Printf("client index=%d, pid=%d get scene successfully\n", c.index, c.ptData.PID)
	return true
}