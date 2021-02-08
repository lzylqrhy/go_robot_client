package net

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/gorilla/websocket"
	"github/go-robot/common"
	"github/go-robot/games"
	"github/go-robot/games/fish"
	"github/go-robot/protocols"
	"io"
	"log"
	"net/url"
	"strings"
	"sync"
)

func WSConnect(ctx context.Context, wg *sync.WaitGroup, sAddr string, pd *common.PlatformData) {
	defer wg.Done()
	log.SetFlags(0)
	u := url.URL{Scheme: "ws", Host: sAddr}
	log.Printf("connect to %s", u.String())
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("web socket dial failed, err:", err)
	}
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
				log.Printf("conn %d eixt\n", index)
				return
			default:
				err := wsReadBuff(conn, chRead)
				if err != nil {
					log.Printf("Read protocol failed, index: %v, err: %v\n", index, err)
					myCancel()
				}
			}
		}
	}()
	// 发数据
	for {
		select {
		case <-myRootCtx.Done():
			conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			log.Printf("conn %d eixt\n", index)
			return
		case pbBuff := <-chWrite:
			if len(pbBuff) >= protocols.HeadSize {
				var buff bytes.Buffer
				buff.Write(pbBuff)
				buff.WriteString("\r\n")
				err = conn.WriteMessage(websocket.TextMessage, buff.Bytes())
				if err != nil {
					log.Println("write failed, err:", err)
					return
				}
			}else {
				myCancel()
			}
		}
	}
}

func wsReadBuff(conn *websocket.Conn, chRead chan<- *protocols.Protocol) error {
	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Println("read failed, err:", err)
		return err
	}
	newMessage := strings.Trim(string(message), "\r\n")
	br := bytes.NewReader([]byte(newMessage))
	for br.Len() > 0 {
		ptData := new(protocols.Protocol)
		err = binary.Read(br, binary.LittleEndian, &ptData.Head)
		if err != nil {
			log.Println("binary.read failed, error is ", err)
			return err
		}
		if ptData.Head.Len > protocols.HeadSize {
			leftBuff := make([]byte, ptData.Head.Len-protocols.HeadSize)
			if _, err := io.ReadFull(br, leftBuff); err != nil {
				fmt.Println("read protocol content failed, error is ", err)
				return err
			}
			ptData.Content.Write(leftBuff)
		}
		chRead <- ptData
	}
	return nil
}

