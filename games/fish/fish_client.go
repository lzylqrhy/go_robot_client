package fish

import (
	"fmt"
	"github/go-robot/common"
	"github/go-robot/protocols"
)

type FClient struct {
	common.ClientBase
}

func NewClient(index uint32, pd *common.PlatformData) common.Client {
	c := new(FClient)
	c.Index = index
	c.PtData = pd
	fmt.Printf("new client: %v \n", c)
	fmt.Printf("new client: %v \n", *c.PtData)
	return c
}

func (c *FClient)Update(ch chan<- []byte) {
	//fmt.Printf("client serial=%d update\n", c.serial)
}

func (c *FClient)OnConnected(ch chan<- []byte)  {
	var ping protocols.C2SSyncTime
	ch<- ping.Bytes()
	fmt.Printf("client index=%d connected\n", c.Index)
}

func (c *FClient)OnDisconnected()  {
	fmt.Printf("client index=%d disconnected\n", c.Index)
}

func (c *FClient)ProcessProtocols(ch chan<- []byte, p *protocols.Protocol) bool {
	fmt.Printf("cmd:0x%04x, data: %v\n", p.Head.Cmd, p.Content)
	isCommon, isOk := c.ProcessCommonProtocols(ch, p)
	if isCommon {
		return isOk
	}
	switch p.Head.Cmd {
	case protocols.PlayerCode:
		return c.processPlayerInfo(ch, p)
	case protocols.EnterRoomCode:
		return c.processEnterRoom(ch, p)
	case protocols.SceneInfoCode:
		return c.processSceneInfo(ch, p)
	case protocols.PlayerSeatCode:
		return c.processSeatsInfo(ch, p)
	case protocols.FishListCode:
		return c.processFishList(ch, p)
	case protocols.BulletListCode:
		return c.processBulletList(ch, p)
	}
	return true
}

func (c *FClient) processPlayerInfo(_ chan<- []byte, p *protocols.Protocol) bool {
	var s2cPlayer protocols.S2CPlayerInfo
	s2cPlayer.Parse(p)
	fmt.Printf("client index=%d, pid=%d get player info successfully, player=%v\n", c.Index, c.PtData.PID, s2cPlayer)
	return true
}

func (c *FClient) processEnterRoom(ch chan<- []byte, p *protocols.Protocol) bool {
	var s2cEnterRoom protocols.S2CEnterRoom
	s2cEnterRoom.Parse(p)
	if s2cEnterRoom.Result != 0 {
		// 进入失败
		fmt.Printf("client index=%d, pid=%d enter room=%d failed, result: %d\n",
			c.Index, c.PtData.PID, s2cEnterRoom.RoomID, s2cEnterRoom.Result)
		return false
	}
	fmt.Printf("client index=%d, pid=%d enter room=%d successfully\n", c.Index, c.PtData.PID, s2cEnterRoom.RoomID)
	// 请求场景信息
	var c2sSceneInfo protocols.C2SGetSceneInfo
	ch<- c2sSceneInfo.Bytes()
	// 转发
	var c2sTransmit protocols.C2STransmitActivity
	c2sTransmit.Activity = "hello"
	ch<- c2sTransmit.Bytes()
	return true
}

func (c *FClient) processSceneInfo(ch chan<- []byte, p *protocols.Protocol) bool {
	var s2cSceneInfo protocols.S2CGetSceneInfo
	s2cSceneInfo.Parse(p)
	fmt.Printf("client index=%d, pid=%d get scene info successfully\n", c.Index, c.PtData.PID)
	return true
}

func (c *FClient) processSeatsInfo(ch chan<- []byte, p *protocols.Protocol) bool {
	var s2cSeats protocols.S2CSeatsInfo
	s2cSeats.Parse(p)
	fmt.Printf("client index=%d, pid=%d get seat list successfully\n", c.Index, c.PtData.PID)
	return true
}

func (c *FClient) processFishList(ch chan<- []byte, p *protocols.Protocol) bool {
	var s2cFish protocols.S2CFishList
	s2cFish.Parse(p)
	fmt.Printf("client index=%d, pid=%d get fish list successfully\n", c.Index, c.PtData.PID)
	fmt.Println(s2cFish.FishList)
	return true
}

func (c *FClient) processBulletList(ch chan<- []byte, p *protocols.Protocol) bool {
	var s2cBullet protocols.S2CBulletList
	s2cBullet.Parse(p)
	fmt.Printf("client index=%d, pid=%d get bullet list successfully\n", c.Index, c.PtData.PID)
	return true
}