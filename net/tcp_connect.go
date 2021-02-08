package net

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/binary"
	"fmt"
	"github/go-robot/common"
	"github/go-robot/games"
	"github/go-robot/games/fish"
	"github/go-robot/protocols"
	"github/go-robot/util"
	"io"
	"io/ioutil"
	"net"
	"sync"
	"time"
)

func TcpConnect(ctx context.Context, wg *sync.WaitGroup, sAddr string, pd *common.PlatformData) {
	defer wg.Done()
	tcpAddr, err := net.ResolveTCPAddr("tcp", sAddr)
	util.CheckError(err)
	conn, err := net.Dial("tcp", tcpAddr.String())
	util.CheckError(err)
	defer conn.Close()
	// 读取上下文数据
	index := ctx.Value("index").(int)
	// 为当前连接创建新的根context
	myRootCtx, myCancel := context.WithCancel(ctx)
	defer myCancel()
	// 创建收发channel
	chRead := make(chan *protocols.Protocol, 100)
	chWrite := make(chan []byte, 10)
	// 创建客户端
	var c games.Client = fish.NewClient(uint32(index), pd)
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
			select {
			case <-myRootCtx.Done():
				fmt.Printf("conn %d eixt\n", index)
				return
			default:
				data, err := readBuff(conn)
				if err != nil {
					fmt.Printf("Read protocol failed, index: %v, err: %v\n", index, err)
					myCancel()
				}
				chRead <- data
			}
		}
	}()
	// 发数据
	for {
		select {
		case <-myRootCtx.Done():
			fmt.Printf("conn %d eixt\n", index)
			return
		case pbBuff := <-chWrite:
			if len(pbBuff) >= protocols.HeadSize {
				_, err = conn.Write(pbBuff)
				if err != nil {
					fmt.Println("write failed, err:", err)
					return
				}
			}else {
				myCancel()
			}
		}
	}
}

func ZipDecode(src []byte) []byte {
	br := bytes.NewReader(src)
	r, err := zlib.NewReaderDict(br, []byte("FK3G"))
	util.CheckError(err)
	defer r.Close()
	dst, err := ioutil.ReadAll(r)
	util.CheckError(err)
	return dst
}

func readBuff(conn net.Conn) (*protocols.Protocol, error) {
	buff := make([]byte, protocols.HeadSize)
	if _, err := io.ReadFull(conn, buff); err != nil {
		fmt.Println("read protocol header failed, error is ", err)
		return nil, err
	}
	ptData := new(protocols.Protocol)
	r := bytes.NewReader(buff)
	err := binary.Read(r, binary.LittleEndian, &ptData.Head)
	if err != nil {
		fmt.Println("binary.read failed, error is ", err)
		return nil, err
	}
	if ptData.Head.Len > protocols.HeadSize {
		leftBuff := make([]byte, ptData.Head.Len - protocols.HeadSize)
		if _, err := io.ReadFull(conn, leftBuff); err != nil {
			fmt.Println("read protocol content failed, error is ", err)
			return nil, err
		}
		//const mark uint16 = 0x8000
		//if ptData.Head.Cmd & mark == mark {
		//	ptData.Head.Cmd &= ^mark
		//	ptData.Content.Write(ZipDecode(leftBuff))
		//}else {
		//	ptData.Content.Write(leftBuff)
		//}
		ptData.Content.Write(leftBuff)
	}
	return ptData, nil
}

func processProtocol(ctx context.Context, wg *sync.WaitGroup, chRead <-chan *protocols.Protocol, chWrite chan<- []byte)  {
	defer wg.Done()
	c, bOk := ctx.Value("player").(games.Client)
	if !bOk {
		fmt.Println("type of context's value is not a client pointer")
		return
	}
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
}