package aladdin

import (
	"github/go-robot/common"
	"github/go-robot/core"
	mynet "github/go-robot/core/mynet"
	"github/go-robot/core/protocol"
	"github/go-robot/global"
	"github/go-robot/global/ini"
	"github/go-robot/protocols"
	"log"
	"math/rand"
	"time"
)

type FClient struct {
	common.ClientBase        // 公共数据
	charID            uint32
	gameCurrency      uint64	// 游戏币
}

func NewClient(index uint, pd *common.PlatformData, dialer mynet.MyDialer) core.RobotClient {
	c := new(FClient)
	c.Index = uint32(index)
	c.PtData = pd
	c.Dialer = dialer
	log.Printf("new client: %v \n", c)
	log.Printf("new client: %v \n", *c.PtData)
	return c
}

// 获取服务端的毫秒
func (c *FClient) getServerTime() uint64 {
	if 0 == c.LocalTime {
		return 0
	}
	return c.SevTime + (uint64(time.Now().UnixNano() / 1e6) - c.LocalTime)
}

func (c *FClient)Update() {
	//log.Printf("client serial=%d update\n", c.serial)
	//if c.IsWorking {
	//	c.play()
	//}
}

func (c *FClient)OnConnected()  {
	var ping protocols.C2SSyncTime
	c.SendPacket(ping.Bytes())
	log.Printf("client index=%d connected\n", c.Index)
}

func (c *FClient)OnDisconnected()  {
	log.Printf("client index=%d disconnected\n", c.Index)
}

func (c *FClient)ProcessProtocols(p *protocol.Protocol) bool {
	//log.Printf("cmd:0x%04x\n", p.Head.Cmd)
	isCommon, isOK := c.ProcessCommonProtocols(p)
	if isCommon {
		return isOK
	}
	switch p.Head.Cmd {
	case protocols.S2CLoginCode:
		return c.processLogin(p)
	case protocols.AladdinPlayerCode:
		return c.processPlayerInfo(p)
	case protocols.EnterRoomCode:
		return c.processEnterRoom(p)
	case protocols.AladdinJoinRoomCode:
		return c.processJoinRoom(p)
	case protocols.AladdinPlayCode:
		return c.processPlayResult(p)
	}
	log.Printf("cmd:0x%04x don't process\n", p.Head.Cmd)
	return true
}

func (c *FClient) processLogin(p *protocol.Protocol) bool {
	var s2cLogin protocols.S2CLogin
	s2cLogin.Parse(p)
	if s2cLogin.Status == 1 {
		// 登录成功
		log.Printf("client index=%d, pid=%d login successfully\n", c.Index, c.PtData.PID)
		// 发送资源加载完成
		var c2sLoaded protocols.C2SResourceLoaded
		c.SendPacket(c2sLoaded.Bytes())
		// 进入房间
		var c2sEnterRoom protocols.C2SEnterRoom
		c2sEnterRoom.RoomID = uint32(ini.AladdinSetting.RoomID)
		c.SendPacket(c2sEnterRoom.Bytes())
		return true
	}
	log.Printf("client index=%d, pid=%d login failed, status: %d\n", c.Index, c.PtData.PID, s2cLogin.Status)
	return false
}

func (c *FClient) processPlayerInfo(p *protocol.Protocol) bool {
	var s2cPlayer protocols.S2CAladdinPlayerInfo
	s2cPlayer.Parse(p)
	log.Printf("client index=%d, pid=%d get player info successfully, player=%v\n", c.Index, c.PtData.PID, s2cPlayer)
	c.charID = s2cPlayer.CharID
	return true
}

func (c *FClient) processEnterRoom(p *protocol.Protocol) bool {
	var s2cEnterRoom protocols.S2CEnterRoom
	s2cEnterRoom.Parse(p)
	if s2cEnterRoom.Result != 1 {
		log.Printf("client index=%d, pid=%d enter room failed\n", c.Index, c.PtData.PID)
	}
	c.RoomID = s2cEnterRoom.RoomID
	// 加入房间
	var c2sJoinRoom protocols.C2SAladdinJoinRoom
	c2sJoinRoom.GameID = 4
	c.SendPacket(c2sJoinRoom.Bytes())
	log.Printf("client index=%d, pid=%d enter room successfully\n", c.Index, c.PtData.PID)
	return true
}

func (c *FClient) processJoinRoom(p *protocol.Protocol) bool {
	var s2cJoinRoom protocols.S2CAladdinJoinRoom
	s2cJoinRoom.Parse(p)
	if s2cJoinRoom.Result == 1 {
		// 加入房间成功
		log.Printf("client index=%d, pid=%d join room successfully\n", c.Index, c.PtData.PID)
		// 拷贝钱
		c.gameCurrency = c.Items[global.ItemCoin]
		// 游戏
		c.play()
		return true
	}
	log.Printf("client index=%d, pid=%d join room failed, status: %d\n", c.Index, c.PtData.PID, s2cJoinRoom.Result)
	return false
}

func (c *FClient) processPlayResult(p *protocol.Protocol) bool {
	var s2cPlay protocols.S2CAladdinPlayResult
	s2cPlay.Parse(p)
	if s2cPlay.Result == 1 {
		// 下注成功，返回游戏结果
		for _, re := range s2cPlay.ReInfo {
			c.gameCurrency += uint64(re.Amount)
		}
		log.Printf("client index=%d, pid=%d play result successfully, money=%d\n", c.Index, c.PtData.PID, c.gameCurrency)
		n := time.Duration(rand.Intn(3) + 4)
		time.Sleep(n * time.Second)
		// 游戏
		c.play()
		return true
	}
	log.Printf("client index=%d, pid=%d play result failed, status: %d\n", c.Index, c.PtData.PID, s2cPlay.Result)
	return false
}

func (c *FClient) play() {
	line, chip := uint8(ini.AladdinSetting.Line), uint32(ini.AladdinSetting.Chip)
	allAmount := uint32(line) * chip
	if uint64(allAmount) > c.gameCurrency {
		log.Printf("client index=%d, pid=%d has no enough money, pull failed, need=%d, cur=%d\n",
			c.Index, c.PtData.PID, allAmount, c.gameCurrency)
		return
	}
	c.gameCurrency -= uint64(allAmount)
	// 发游戏协议
	c2sPlay := protocols.C2SAladdinPlay{}
	c2sPlay.Line = uint8(ini.AladdinSetting.Line)
	c2sPlay.Amount = allAmount
	c.SendPacket(c2sPlay.Bytes())
}

