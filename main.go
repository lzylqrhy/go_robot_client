package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github/go-robot/protocols"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
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

func createConnect(ctx context.Context, wg *sync.WaitGroup, sAddr string) {
	defer wg.Done()
	conn, err := net.Dial("tcp", sAddr)
	if err != nil {
		fmt.Println("connect server failed, addr: ", sAddr)
		return
	}
	defer conn.Close()
	// 读取上下文数据
	index, _ := ctx.Value("index").(int)
	// 处理协议
	wg.Add(1)
	subCtx := context.WithValue(ctx, "player", NewClient(index))
	ch := make(chan []byte, 10)
	go ProcessProtocol(subCtx, wg, ch)
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("conn %d eixt", index)
			return
		case pbBuff := <-ch:
			fmt.Println(pbBuff)
			// 发消息
			_, err = conn.Write(pbBuff)
			if err != nil {
				fmt.Println("write failed, err:", err)
				return
			}
		default:
			// 读消息
			data, err := readBuff(conn)
			if err != nil {
				fmt.Printf("Read protocol failed, index: %v, err: %v\n", index, err)
				break
			}
			ch<- data
		}
	}
}

func readBuff(conn net.Conn) ([]byte, error) {
	var data bytes.Buffer
	buff := make([]byte,4)
	if _, err := io.ReadFull(conn, buff); err != nil {
		fmt.Println("read protocol header failed, error is ", err)
		return data.Bytes(), err
	}
	var header protocols.ProtocolHeader
	r := bytes.NewReader(buff)
	err := binary.Read(r, binary.LittleEndian, &header)
	if err != nil {
		fmt.Println("binary.read failed, error is ", err)
		return data.Bytes(), err
	}
	data.Write(buff)
	if header.Len > 4 {
		leftBuff := make([]byte, header.Len - 4)
		if _, err := io.ReadFull(conn, leftBuff); err != nil {
			fmt.Println("read protocol content failed, error is ", err)
			return data.Bytes(), err
		}
		data.Write(leftBuff)
	}
	return data.Bytes(), nil
}

func ProcessProtocol(ctx context.Context, wg *sync.WaitGroup, ch chan []byte)  {
	defer wg.Done()
	c := ctx.Value("player").(*Client)
	frame := time.NewTicker(time.Millisecond * 200)
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("process protocol %d eixt", c.serial)
			return
		case pbBuff := <-ch:
			c.ProcessProtocols(ch, pbBuff)
		case <-frame.C:
			c.Update()
		}
	}
}