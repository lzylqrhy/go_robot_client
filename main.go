package main

import (
	"context"
	"github/go-robot/common"
	"github/go-robot/core"
	"github/go-robot/core/mynet"
	"github/go-robot/games"
	"github/go-robot/global"
	"github/go-robot/global/ini"
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
	// 设置全局随机种子
	rand.Seed(time.Now().Unix())
	// 设置log文件
	f, err := os.OpenFile("run.log", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)
	defer f.Close()

	// 读取配置
	ini.LoadSetting()
	cfg := &ini.MainSetting
	// 从平台获取信息
	userList := common.GetPlatformUserData()
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	// 设置机器人数据
	games.SetRobotTestData(userList)
	// 启动机器人
	for i, user := range userList {
		// 连接服务器
		var serverAddr string
		switch cfg.GameID {
		case global.FishGame:
			serverAddr = ini.FishSetting.ServerAddr
		case global.FruitGame:
			serverAddr = ini.FruitSetting.ServerAddr
		}
		if "" == serverAddr {
			if mynet.WS == cfg.NetProtocol {
				serverAddr = user.WSServerAddr
			} else {
				log.Println("when protocol is not ws, must set server")
				break
			}
		}
		// 网络连接
		d := mynet.NewConnect(cfg.NetProtocol, serverAddr)
		// 创建客户端
		c := games.NewClient(uint(i), user, d)
		// 开工
		core.DoWork(ctx, &wg, c, d, ini.GameCommonSetting.Frame)
	}
	//http.ListenAndServe("0.0.0.0:6060", nil)
	// 监听信号
	go func() {
		waitForASignal()
		cancel()
		log.Println("stop all jobs")
	}()
	wg.Wait()
	log.Println("exit")
}
// 等待主动关闭的信号
func waitForASignal()  {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)
	<-sig
}

