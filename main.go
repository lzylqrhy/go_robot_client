package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	// 读取配置
	servAddr := "192.168.0.194:7710"
	// 机器人数量
	playerNum := 1
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	for i := 0; i < playerNum; i++ {
		wg.Add(1)
		subCtx := context.WithValue(ctx, "index", i)
		go createConnect(subCtx, &wg, servAddr)
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
