package main

import (
	"context"
	"github/go-robot/common"
	"github/go-robot/games"
	myNet "github/go-robot/net"
	"github/go-robot/util"
	"gopkg.in/ini.v1"
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
	conf, err := ini.Load("./configs/main.ini")
	util.CheckError(err)
	netProtocol := conf.Section("server").Key("protocol").String()
	serverAddr := conf.Section("server").Key("server_addr").String()
	// 机器人数量
	userStart, err := conf.Section("robot").Key("start").Uint()
	util.CheckError(err)
	num, err := conf.Section("robot").Key("num").Uint()
	util.CheckError(err)
	gameID, err := conf.Section("robot").Key("game_id").Uint()
	util.CheckError(err)

	// 从平台获取信息
	userList := common.GetPlatformUserData(userStart, num)
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	for i, user := range userList {
		// 连接服务器
		d := myNet.NewConnect(netProtocol, serverAddr)
		// 创建客户端
		c := games.NewClient(gameID, uint(i), user, d)
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

