package main

import (
	"context"
	"github/go-robot/common"
	"github/go-robot/games"
	"github/go-robot/global"
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
	global.LoadSetting()
	cfg := &global.MainSetting
	// 从平台获取信息
	userList := common.GetPlatformUserData()
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	for i, user := range userList {
		// 连接服务器
		var serverAddr string
		switch cfg.GameID {
		case games.FishGame:
			serverAddr = global.FishSetting.ServerAddr
		}
		if "" == serverAddr {
			serverAddr = user.ServerAddr
		}
		d := myNet.NewConnect(cfg.NetProtocol, serverAddr)
		// 创建客户端
		c := games.NewClient(cfg.GameID, uint(i), user, d)
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

