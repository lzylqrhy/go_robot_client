/**
 机器人客户端接口，游戏主逻辑入口
 created by lzy
*/
package core

import (
	"context"
	"github/go-robot/core/mynet"
	"github/go-robot/core/protocol"
	"log"
	"runtime/debug"
	"sync"
	"time"
)

// 机器人客户端接口
type RobotClient interface {
	Update()
	OnConnected()
	OnDisconnected()
	ProcessProtocols(p *protocol.Protocol) bool
}

// 游戏主入口
func DoWork(ctx context.Context, wg *sync.WaitGroup, c RobotClient, d mynet.MyDialer, frameCount uint)  {
	myCtx, cancel := context.WithCancel(ctx)
	wg.Add(1)
	go func() {
		defer func() {
			cancel()
			wg.Done()
			if r := recover(); nil != r {
				log.Println("run time error was caught: ", r)
				log.Printf("error stack: %v \n", string(debug.Stack()))
			}
		}()
		// 驱动网络连接
		if !d.Run(myCtx, wg) {
			return
		}
		frameDuration := 1000
		if frameCount > 0 && frameCount <= 1000 {
			frameDuration = int(1000 / frameCount)
		}
		frameTick := time.NewTicker(time.Millisecond * time.Duration(frameDuration))
		pingTick := time.NewTicker(time.Second * 10)
		for {
			select {
			case <-ctx.Done():
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
							// 断线重连，5秒一次
						Break:
							for i := 0; i < 5; i++ {
								select {
								case <-time.After(5 * time.Second):
									log.Println("client reconnect count i = ", i)
									if d.Run(myCtx, wg) {
										break Break
									}
									if 4 == i {
										// 放弃重连，关闭
										return
									}
								}
							}
						}
						break
					default:
						if !c.ProcessProtocols(pd) {
							// 如果处理失败，主动断开连接
							return
						}
					}
				}else {
					// 当前用户退出游戏
					return
				}
			case <-frameTick.C: // 客户端定时器
				c.Update()
			case <-pingTick.C:
				var ping protocol.C2SPing
				ping.TimeStamp = uint32(time.Now().Unix())
				d.SendPacket(ping.Bytes())
			}
		}
	}()
}
