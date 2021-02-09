package common

import (
	"context"
	"fmt"
	"github/go-robot/net"
	"github/go-robot/protocols"
	"strings"
	"sync"
	"time"
)

type Client interface {
	Update(ch chan<- []byte)
	OnConnected(ch chan<- []byte)
	OnDisconnected()
	ProcessProtocols(ch chan<- []byte, p *protocols.Protocol) bool
}

type NClient interface {
	Update(ch chan<- []byte)
	OnConnected(ch chan<- []byte)
	OnDisconnected()
	ProcessProtocols(ch chan<- []byte, p *protocols.Protocol) bool
}

// 客户端基类
type ClientBase struct {
	Index   uint32
	PtData  *PlatformData
	SevTime uint32
	Dialer net.MyDialer
}

// 处理公共协议
func (c *ClientBase)ProcessCommonProtocols(ch chan<- []byte, p *protocols.Protocol) (bool, bool) {
	switch p.Head.Cmd {
	case protocols.SyncTimeCode:
		return true, c.processSyncTime(ch, p)
	case protocols.S2CLoginCode:
		return true, c.processLogin(ch, p)
	}
	return false, false
}

func (c *ClientBase) processSyncTime(ch chan<- []byte, p *protocols.Protocol) bool {
	s2cSync := new(protocols.S2CSyncTime)
	s2cSync.Parse(p)
	fmt.Println(s2cSync)
	c.SevTime = s2cSync.TimeStamp
	// 请求登录
	var s2cLogin protocols.C2SLogin
	s2cLogin.IsChildGame = false
	var strBuilder strings.Builder
	strBuilder.WriteString(c.PtData.LoginToken)
	strBuilder.WriteString(":0x20:1")
	s2cLogin.Token = strBuilder.String()
	fmt.Println("session:", s2cLogin.Token)
	ch<- s2cLogin.Bytes()
	return true
}

func (c *ClientBase) processLogin(ch chan<- []byte, p *protocols.Protocol) bool {
	var s2cLogin protocols.S2CLogin
	s2cLogin.Parse(p)
	if s2cLogin.Status == 1 {
		// 登录成功
		fmt.Printf("client index=%d, pid=%d login successfully\n", c.Index, c.PtData.PID)
		// 发送资源加载完成
		var c2sLoaded protocols.C2SResourceLoaded
		ch<- c2sLoaded.Bytes()
		// 进入房间
		var c2sEnterRoom protocols.C2SEnterRoom
		c2sEnterRoom.RoomID = 20001
		ch<- c2sEnterRoom.Bytes()
		return true
	}
	fmt.Printf("client index=%d, pid=%d login failed, status: %d\n", c.Index, c.PtData.PID, s2cLogin.Status)
	return false
}

func DoWork(ctx context.Context, wg *sync.WaitGroup)  {
	myCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	// 驱动网络连接
	c.Dialer.Run(myCtx, wg)
	wg.Add(1)
	go func() {
		defer wg.Done()

		// 连接成功事件
		c.OnConnected(chWrite)
		frameTick := time.NewTicker(time.Millisecond * 200)
		pingTick := time.NewTicker(time.Second * 10)
		for {
			select {
			case <-ctx.Done(): // 模拟断线事件
				frameTick.Stop()
				c.OnDisconnected()
				return
			case pbBuff := <-chRead: // 处理数据
				if pbBuff != nil {
					switch pbBuff.Head.Cmd {
					case 2:
						break
					default:
						if !c.ProcessProtocols(chWrite, pbBuff) {
							chWrite<- []byte("0")
						}
					}
				}
			case <-frameTick.C: // 客户端定时器
				c.Update(chWrite)
			case <-pingTick.C:
				var ping protocols.C2SPing
				ping.TimeStamp = uint32(time.Now().Unix())
				chWrite<- ping.Bytes()
			}
		}
	}()
}