package main

import (
	"context"
	"github/go-robot/common"
	"github/go-robot/games/fish"
	myNet "github/go-robot/net"
	"log"
	"math/rand"
	//_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	// 读取配置
	// serverAddr := "192.168.0.194:7710"
	//serverAddr := "127.0.0.1:8080" //"192.168.0.194:7712"
	serverAddr := "192.168.0.194:7712"
	// 机器人数量
	userStart := uint(2)
	userEnd := uint(5)
	// 从平台获取信息
	userList := common.GetPlatformUserData(userStart, userEnd)
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	for i, user := range userList {
		//wg.Add(1)
		//subCtx := context.WithValue(ctx, "index", i)
		//go myNet.TcpConnect(subCtx, &wg, serverAddr, user)
		//go myNet.WSConnect(subCtx, &wg, serverAddr, user)

		// 连接服务器
		d := myNet.NewConnect(serverAddr)
		// 创建客户端
		c := fish.NewClient(uint32(i), user, d)
		// 开工
		common.DoWork(ctx, &wg, c, d)
	}
	//http.ListenAndServe("0.0.0.0:6060", nil)
	// 监听信号
	waitForASignal()
	cancel()
	log.Println("stop all jobs")
	wg.Wait()
	log.Println("exit")
}

func waitForASignal()  {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)
	<-sig
}

