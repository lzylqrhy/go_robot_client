package common

import (
	"context"
	myNet "github/go-robot/net"
	"github/go-robot/protocols"
	"log"
	"strings"
	"sync"
	"time"
)

//type Client interface {
//	Update(ch chan<- []byte)
//	OnConnected(ch chan<- []byte)
//	OnDisconnected()
//	ProcessProtocols(ch chan<- []byte, p *protocols.Protocol) bool
//}

type Client interface {
	Update()
	OnConnected()
	OnDisconnected()
	ProcessProtocols(p *protocols.Protocol) bool
}

// 客户端基类
type ClientBase struct {
	Index     uint32
	PtData    *PlatformData
	SevTime   uint64
	LocalTime uint64
	Dialer    myNet.MyDialer
}

// 处理公共协议
func (c *ClientBase)ProcessCommonProtocols(p *protocols.Protocol) (bool, bool) {
	switch p.Head.Cmd {
	case protocols.SyncTimeCode:
		return true, c.processSyncTime(p)
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

func DoWork(ctx context.Context, wg *sync.WaitGroup, c Client, d myNet.MyDialer)  {
	myCtx, cancel := context.WithCancel(ctx)
	wg.Add(1)
	go func() {
		defer func() {
			cancel()
			wg.Done()
		}()
		// 驱动网络连接
		d.Run(myCtx, wg)
		frameTick := time.NewTicker(time.Millisecond * 200)
		pingTick := time.NewTicker(time.Second * 10)
		for {
			select {
			case <-ctx.Done(): // 模拟断线事件
				frameTick.Stop()
				pingTick.Stop()
				return
			case pd := <-d.ReadPacket(): // 处理数据
				if pd != nil {
					switch pd.Head.Cmd {
					case 0:
						switch pd.Head.Len {
						case 0: // 连接成功
							c.OnConnected()
						case 1: // 断开连接
							c.OnDisconnected()
							d.Run(myCtx, wg)
						}
						break
					default:
						if !c.ProcessProtocols(pd) {
							d.Disconnect()
						}
					}
				}else {
					return
				}
			case <-frameTick.C: // 客户端定时器
				c.Update()
			case <-pingTick.C:
				var ping protocols.C2SPing
				ping.TimeStamp = uint32(time.Now().Unix())
				d.SendPacket(ping.Bytes())
			}
		}
	}()
}