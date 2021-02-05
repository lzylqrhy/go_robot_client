package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github/go-robot/protocols"
	"io"
	"net"
	"sync"
	"time"
)

func createConnect(ctx context.Context, wg *sync.WaitGroup, sAddr string) {
	defer wg.Done()
	tcpAddr, err := net.ResolveTCPAddr("tcp", sAddr)
	if err != nil {
		fmt.Println("net.ResolveTCPAddr failed, err: ", err)
		return
	}
	conn, err := net.Dial("tcp", tcpAddr.String())
	if err != nil {
		fmt.Println("connect server failed, addr: ", sAddr)
		return
	}
	defer conn.Close()
	// 读取上下文数据
	index, bOk := ctx.Value("index").(int)
	if !bOk {
		fmt.Println("type of context's value is not int")
		return
	}
	// 为当前连接创建新的根context
	myRootCtx, myCancel := context.WithCancel(ctx)
	defer myCancel()
	// 创建收发channel
	chRead := make(chan []byte, 100)
	chWrite := make(chan []byte, 10)
	// 创建客户端
	c := NewClient(index)
	// 处理协议
	wg.Add(1)
	subCtx := context.WithValue(myRootCtx, "player", c)
	go processProtocol(subCtx, wg, chRead, chWrite)
	// 收数据
	wg.Add(1)
	go func() {
		defer wg.Done()
		// 读消息
		for {
			data, err := readBuff(conn)
			if err != nil {
				fmt.Printf("Read protocol failed, index: %v, err: %v\n", index, err)
				return
			}
			chRead <- data
		}
	}()
	// 发数据
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("conn %d eixt\n", index)
			return
		case pbBuff := <-chWrite:
			_, err = conn.Write(pbBuff)
			if err != nil {
				fmt.Println("write failed, err:", err)
				return
			}
		}
	}
}

func readBuff(conn net.Conn) ([]byte, error) {
	var data bytes.Buffer
	const HeadSize = 4
	buff := make([]byte, HeadSize)
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
	if header.Len > HeadSize {
		leftBuff := make([]byte, header.Len - HeadSize)
		if _, err := io.ReadFull(conn, leftBuff); err != nil {
			fmt.Println("read protocol content failed, error is ", err)
			return data.Bytes(), err
		}
		data.Write(leftBuff)
	}
	return data.Bytes(), nil
}

func processProtocol(ctx context.Context, wg *sync.WaitGroup, chRead <-chan []byte, chWrite chan<- []byte)  {
	defer wg.Done()
	c, bOk := ctx.Value("player").(*Client)
	if !bOk {
		fmt.Println("type of context's value is not a client pointer")
		return
	}
	// 连接成功事件
	c.OnConnected(chWrite)
	frameTick := time.NewTicker(time.Millisecond * 200)
	for {
		select {
		case <-ctx.Done(): // 模拟断线事件
			frameTick.Stop()
			c.OnDisconnected()
			fmt.Printf("process protocol %d eixt\n", c.serial)
			return
		case pbBuff := <-chRead: // 处理数据
			c.ProcessProtocols(chWrite, pbBuff)
		case <-frameTick.C: // 客户端定时器
			c.Update(chWrite)
		}
	}
}