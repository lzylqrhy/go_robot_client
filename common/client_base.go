package common

import (
	myNet "github/go-robot/net"
	"github/go-robot/protocols"
	"log"
	"strings"
)

// 客户端基类
type ClientBase struct {
	Index     uint32	// 机器人索引
	PtData    *PlatformData	// 平台数据
	SevTime   uint64	// 服务端时间
	LocalTime uint64	// 本地时间
	Dialer    myNet.MyDialer	// 网络连接器
	Items map[uint32]uint64	// 物品
	IsWorking bool	// 是否在工作
}

// 处理公共协议
func (c *ClientBase)ProcessCommonProtocols(p *protocols.Protocol) (bool, bool) {
	switch p.Head.Cmd {
	case protocols.SyncTimeCode:
		return true, c.processSyncTime(p)
	case protocols.OpenPackageCode:
		return true, c.processOpenPackage(p)
	}
	return false, false
}

func (c *ClientBase) processSyncTime(p *protocols.Protocol) bool {
	s2cSync := new(protocols.S2CSyncTime)
	s2cSync.Parse(p)
	log.Println(s2cSync)
	// 请求登录
	var s2cLogin protocols.C2SLogin
	s2cLogin.IsChildGame = false
	var strBuilder strings.Builder
	strBuilder.WriteString(c.PtData.LoginToken)
	strBuilder.WriteString(":0x20:1")
	s2cLogin.Token = strBuilder.String()
	log.Println("session:", s2cLogin.Token)
	c.SendPacket(s2cLogin.Bytes())
	return true
}

func (c *ClientBase) SendPacket(msg []byte)  {
	c.Dialer.SendPacket(msg)
}

func (c *ClientBase) processOpenPackage(p *protocols.Protocol) bool {
	var s2cPackage protocols.S2COpenPackage
	s2cPackage.Parse(p)
	if c.Items == nil {
		c.Items = make(map[uint32]uint64)
	}
	for _, item := range s2cPackage.Items {
		c.Items[item.ModeID] = uint64(item.Amount)
	}
	log.Printf("client index=%d, pid=%d package\n", c.Index, c.PtData.PID)
	return true
}

