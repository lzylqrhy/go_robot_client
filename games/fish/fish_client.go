package fish

import (
	"fmt"
	"github/go-robot/common"
	myNet "github/go-robot/net"
	"github/go-robot/protocols"
)

type FClient struct {
	common.ClientBase
}

func NewClient(index uint32, pd *common.PlatformData, dialer myNet.MyDialer) common.Client {
	c := new(FClient)
	c.Index = index
	c.PtData = pd
	c.Dialer = dialer
	fmt.Printf("new client: %v \n", c)
	fmt.Printf("new client: %v \n", *c.PtData)
	return c
}

//func (c *FClient)Update(ch chan<- []byte) {
//	//fmt.Printf("client serial=%d update\n", c.serial)
//}
//
//func (c *FClient)OnConnected(ch chan<- []byte)  {
//	var ping protocols.C2SSyncTime
//	ch<- ping.Bytes()
//	fmt.Printf("client index=%d connected\n", c.Index)
//}
//
//func (c *FClient)OnDisconnected()  {
//	fmt.Printf("client index=%d disconnected\n", c.Index)
//}
//
//func (c *FClient)ProcessProtocols(ch chan<- []byte, p *protocols.Protocol) bool {
//	fmt.Printf("cmd:0x%04x, data: %v\n", p.Head.Cmd, p.Content)
//	isCommon, isOk := c.ProcessCommonProtocols(ch, p)
//	if isCommon {
//		return isOk
//	}
//	switch p.Head.Cmd {
//	case protocols.PlayerCode:
//		return c.processPlayerInfo(ch, p)
//	case protocols.EnterRoomCode:
//		return c.processEnterRoom(ch, p)
//	case protocols.SceneInfoCode:
//		return c.processSceneInfo(ch, p)
//	case protocols.PlayerSeatCode:
//		return c.processSeatsInfo(ch, p)
//	case protocols.FishListCode:
//		return c.processFishList(ch, p)
//	case protocols.BulletListCode:
//		return c.processBulletList(ch, p)
//	}
//	return true
//}
//
//func (c *FClient) processPlayerInfo(_ chan<- []byte, p *protocols.Protocol) bool {
//	var s2cPlayer protocols.S2CPlayerInfo
//	s2cPlayer.Parse(p)
//	fmt.Printf("client index=%d, pid=%d get player info successfully, player=%v\n", c.Index, c.PtData.PID, s2cPlayer)
//	return true
//}
//
//func (c *FClient) processEnterRoom(ch chan<- []byte, p *protocols.Protocol) bool {
//	var s2cEnterRoom protocols.S2CEnterRoom
//	s2cEnterRoom.Parse(p)
//	if s2cEnterRoom.Result != 0 {
//		// 进入失败
//		fmt.Printf("client index=%d, pid=%d enter room=%d failed, result: %d\n",
//			c.Index, c.PtData.PID, s2cEnterRoom.RoomID, s2cEnterRoom.Result)
//		return false
//	}
//	fmt.Printf("client index=%d, pid=%d enter room=%d successfully\n", c.Index, c.PtData.PID, s2cEnterRoom.RoomID)
//	// 请求场景信息
//	var c2sSceneInfo protocols.C2SGetSceneInfo
//	ch<- c2sSceneInfo.Bytes()
//	// 转发
//	var c2sTransmit protocols.C2STransmitActivity
//	c2sTransmit.Activity = "hello"
//	ch<- c2sTransmit.Bytes()
//	return true
//}
//
//func (c *FClient) processSceneInfo(ch chan<- []byte, p *protocols.Protocol) bool {
//	var s2cSceneInfo protocols.S2CGetSceneInfo
//	s2cSceneInfo.Parse(p)
//	fmt.Printf("client index=%d, pid=%d get scene info successfully\n", c.Index, c.PtData.PID)
//	return true
//}
//
//func (c *FClient) processSeatsInfo(ch chan<- []byte, p *protocols.Protocol) bool {
//	var s2cSeats protocols.S2CSeatsInfo
//	s2cSeats.Parse(p)
//	fmt.Printf("client index=%d, pid=%d get seat list successfully\n", c.Index, c.PtData.PID)
//	return true
//}
//
//func (c *FClient) processFishList(ch chan<- []byte, p *protocols.Protocol) bool {
//	var s2cFish protocols.S2CFishList
//	s2cFish.Parse(p)
//	fmt.Printf("client index=%d, pid=%d get fish list successfully\n", c.Index, c.PtData.PID)
//	fmt.Println(s2cFish.FishList)
//	return true
//}
//
//func (c *FClient) processBulletList(ch chan<- []byte, p *protocols.Protocol) bool {
//	var s2cBullet protocols.S2CBulletList
//	s2cBullet.Parse(p)
//	fmt.Printf("client index=%d, pid=%d get bullet list successfully\n", c.Index, c.PtData.PID)
//	return true
//}

func (c *FClient)Update() {
	//fmt.Printf("client serial=%d update\n", c.serial)
}

func (c *FClient)OnConnected()  {
	var ping protocols.C2SSyncTime
	c.SendPacket(ping.Bytes())
	fmt.Printf("client index=%d connected\n", c.Index)
}

func (c *FClient)OnDisconnected()  {
	fmt.Printf("client index=%d disconnected\n", c.Index)
}

func (c *FClient)ProcessProtocols(p *protocols.Protocol) bool {
	fmt.Printf("cmd:0x%04x, data: %v\n", p.Head.Cmd, p.Content)
	isCommon, isOk := c.ProcessCommonProtocols(p)
	if isCommon {
		return isOk
	}
	switch p.Head.Cmd {
	case protocols.S2CLoginCode:
		return c.processLogin(p)
	case protocols.PlayerCode:
		return c.processPlayerInfo(p)
	case protocols.EnterRoomCode:
		return c.processEnterRoom(p)
	case protocols.SceneInfoCode:
		return c.processSceneInfo(p)
	case protocols.PlayerSeatCode:
		return c.processSeatsInfo(p)
	case protocols.FishListCode:
		return c.processFishList(p)
	case protocols.BulletListCode:
		return c.processBulletList(p)
	}
	return true
}

func (c *FClient) processLogin(p *protocols.Protocol) bool {
	var s2cLogin protocols.S2CLogin
	s2cLogin.Parse(p)
	if s2cLogin.Status == 1 {
		// 登录成功
		fmt.Printf("client index=%d, pid=%d login successfully\n", c.Index, c.PtData.PID)
		// 发送资源加载完成
		var c2sLoaded protocols.C2SResourceLoaded
		c.SendPacket(c2sLoaded.Bytes())
		// 进入房间
		var c2sEnterRoom protocols.C2SEnterRoom
		c2sEnterRoom.RoomID = 20001
		c.SendPacket(c2sEnterRoom.Bytes())
		return true
	}
	fmt.Printf("client index=%d, pid=%d login failed, status: %d\n", c.Index, c.PtData.PID, s2cLogin.Status)
	return false
}

func (c *FClient) processPlayerInfo(p *protocols.Protocol) bool {
	var s2cPlayer protocols.S2CPlayerInfo
	s2cPlayer.Parse(p)
	fmt.Printf("client index=%d, pid=%d get player info successfully, player=%v\n", c.Index, c.PtData.PID, s2cPlayer)
	return true
}

func (c *FClient) processEnterRoom(p *protocols.Protocol) bool {
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
	c.SendPacket(c2sSceneInfo.Bytes())
	//// 转发
	//var c2sTransmit protocols.C2STransmitActivity
	//c2sTransmit.Activity = "hello"
	//c.sendPacket(c2sTransmit.Bytes())
	return true
}

func (c *FClient) processSceneInfo(p *protocols.Protocol) bool {
	var s2cSceneInfo protocols.S2CGetSceneInfo
	s2cSceneInfo.Parse(p)
	fmt.Printf("client index=%d, pid=%d get scene info successfully\n", c.Index, c.PtData.PID)
	return true
}

func (c *FClient) processSeatsInfo(p *protocols.Protocol) bool {
	var s2cSeats protocols.S2CSeatsInfo
	s2cSeats.Parse(p)
	fmt.Printf("client index=%d, pid=%d get seat list successfully\n", c.Index, c.PtData.PID)
	return true
}

func (c *FClient) processFishList(p *protocols.Protocol) bool {
	var s2cFish protocols.S2CFishList
	s2cFish.Parse(p)
	fmt.Printf("client index=%d, pid=%d get fish list successfully\n", c.Index, c.PtData.PID)
	fmt.Println(s2cFish.FishList)
	return true
}

func (c *FClient) processBulletList(p *protocols.Protocol) bool {
	var s2cBullet protocols.S2CBulletList
	s2cBullet.Parse(p)
	fmt.Printf("client index=%d, pid=%d get bullet list successfully\n", c.Index, c.PtData.PID)
	return true
}