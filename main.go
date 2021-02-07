package main

import (
	"context"
	"fmt"
	myNet "github/go-robot/net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	// 读取配置
	servAddr := "192.168.0.194:7710"
	// 机器人数量
	userStart := uint(2)
	userEnd := uint(2)
	// 从平台获取信息
	userList := myNet.GetPlatformUserData(userStart, userEnd)
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	for i, user := range userList {
		wg.Add(1)
		subCtx := context.WithValue(ctx, "index", i)
		go myNet.TcpConnect(subCtx, &wg, servAddr, user)
	}
	// 监听信号
	waitForASignal()
	cancel()
	fmt.Println("stop all jobs")
	wg.Wait()
	fmt.Println("exit")
}

func waitForASignal()  {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)
	<-sig
}

