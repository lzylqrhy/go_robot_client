package net

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/gorilla/websocket"
	"github/go-robot/protocols"
	"io"
	"log"
	"net/url"
	"sync"
)

const postfix = "\r\n\r\n"

type WSDialer struct {
	conn *websocket.Conn
	chRead chan *protocols.Protocol
	chWrite chan []byte
	sAddr string
	ctx context.Context
}

func NewConnect(sAddr string) MyDialer {
	d := new(WSDialer)
	u := url.URL{Scheme: "ws", Host: sAddr}
	d.sAddr = u.String()
	d.chRead = make(chan *protocols.Protocol, 100)
	d.chWrite = make(chan []byte, 100)
	return d
}

func (d *WSDialer) connect() bool {
	var err error
	d.conn, _, err = websocket.DefaultDialer.Dial(d.sAddr, nil)
	if err != nil {
		log.Fatal("web socket dial failed, err:", err)
		return false
	}
	return true
}

func (d *WSDialer) Disconnect() {
	if d.conn != nil {
		err := d.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Println("web socket disconnect failed, err:", err)
		}
		select {
		case <-d.chRead:
			close(d.chWrite)
			close(d.chRead)
		}
	}
}

func (d *WSDialer) close() {
	if d.conn != nil {
		err := d.conn.Close()
		if err != nil {
			log.Println("web socket close socket failed, err:", err)
		}
		d.conn = nil
	}
}

func (d *WSDialer) SendPacket(data []byte) bool {
	if d.conn == nil {
		log.Panicln("send failed, the dialer is offline")
		return false
	}
	d.chWrite<- data
	return true
}

func (d *WSDialer) ReadPacket() <-chan *protocols.Protocol {
	return d.chRead
}

func (d *WSDialer) Run(ctx context.Context, wg *sync.WaitGroup) bool {
	// 先连接
	if !d.connect() {
		return false
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		// 连接成功
		var pd protocols.Protocol
		pd.Head.Cmd = 0
		pd.Head.Len = 0
		d.chRead <- &pd
		// 读消息
		go func() {
			defer d.close()
			d.read()
			// 断开连接
			fmt.Println("send on disconnected message")
			var pd protocols.Protocol
			pd.Head.Cmd = 0
			pd.Head.Len = 1
			d.chRead <- &pd
		}()
		// 发数据
		for {
			select {
			case <-ctx.Done():
				d.Disconnect()
				fmt.Println("conn exit")
				return
			case data := <-d.chWrite:
				if data == nil {
					return
				}
				//enMessage := base64.StdEncoding.EncodeToString(data)
				//enMessage += postfix
				//err := d.conn.WriteMessage(websocket.TextMessage, []byte(enMessage))
				err := d.conn.WriteMessage(websocket.BinaryMessage, data)
				if err != nil {
					fmt.Println("write failed, err:", err)
					return
				}
			}
		}
	}()
	return true
}

func (d *WSDialer) read() {
	// 读消息
	var ptData *protocols.Protocol
	headBuff := make([]byte, 0, protocols.HeadSize)
	leftCount := uint16(0)
	for {
		_, message, err := d.conn.ReadMessage()
		if err != nil {
			fmt.Println("read failed, err:", err)
			break
		}
		//// 去掉
		//newMessage := strings.Trim(string(message), postfix)
		//// 解码
		//deBuff, err:= base64.StdEncoding.DecodeString(newMessage)
		//if err != nil {
		//	log.Println("base64 decoding failed, err:", err)
		//	d.Disconnect()
		//	continue
		//}
		//br := bytes.NewReader(deBuff)
		br := bytes.NewReader(message)
		for {
			if nil == ptData {
				if br.Len() + int(leftCount) < protocols.HeadSize {
					n, err := br.Read(headBuff[leftCount:])
					if err != nil {
						fmt.Println("read head buff failed, error is ", err)
						return
					}
					headBuff = headBuff[:leftCount+uint16(n)]
					leftCount = protocols.HeadSize - (leftCount+uint16(n))
					break
				}
				ptData = new(protocols.Protocol)
				err = binary.Read(br, binary.LittleEndian, &ptData.Head)
				if err != nil {
					fmt.Println("binary.read failed, error is ", err)
					return
				}
				if ptData.Head.Len > protocols.HeadSize {
					leftBuff := make([]byte, ptData.Head.Len-protocols.HeadSize)
					if n, err := io.ReadFull(br, leftBuff); err != nil {
						if err == io.ErrUnexpectedEOF {
							leftCount = ptData.Head.Len - protocols.HeadSize - uint16(n)
							ptData.Content.Write(leftBuff)
							break
						}else {
							fmt.Println("read protocol content failed, error is ", err)
							return
						}
					}
					ptData.Content.Write(leftBuff)
				}
				d.chRead <- ptData
				ptData = nil
			}else {
				if br.Len() < int(leftCount) {
					n, err := br.Read(headBuff[leftCount:])
					if err != nil {
						fmt.Println("read head buff failed, error is ", err)
						return
					}
					headBuff = headBuff[:leftCount+uint16(n)]
					leftCount = protocols.HeadSize - (leftCount+uint16(n))
					break
				}
			}


			if br.Len() > 0 {
				ptData := new(protocols.Protocol)
				err = binary.Read(br, binary.LittleEndian, &ptData.Head)
				if err != nil {
					fmt.Println("binary.read failed, error is ", err)
					return
				}
				if ptData.Head.Len > protocols.HeadSize {
					leftBuff := make([]byte, ptData.Head.Len-protocols.HeadSize)
					if _, err := io.ReadFull(br, leftBuff); err != nil {
						fmt.Println("read protocol content failed, error is ", err)
						return
					}
					ptData.Content.Write(leftBuff)
				}
				d.chRead <- ptData
			} else {
				break
			}
		}
	}
}

//func WSConnect(ctx context.Context, wg *sync.WaitGroup, sAddr string, pd *common.PlatformData) {
//	defer wg.Done()
//	log.SetFlags(0)
//	u := url.URL{Scheme: "ws", Host: sAddr}
//	//u := url.URL{Scheme: "ws", Host: sAddr, Path: "/echo"}
//	log.Printf("connect to %s", u.String())
//	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
//	if err != nil {
//		log.Fatal("web socket dial failed, err:", err)
//	}
//	defer conn.Close()
//	// 读取上下文数据
//	index := ctx.Value("index").(int)
//	// 为当前连接创建新的根context
//	myRootCtx, myCancel := context.WithCancel(ctx)
//	defer myCancel()
//	// 创建收发channel
//	chRead := make(chan *protocols.Protocol, 100)
//	chWrite := make(chan []byte, 10)
//	// 创建客户端
//	c := fish.NewClient(uint32(index), pd)
//	// 处理协议
//	wg.Add(1)
//	subCtx := context.WithValue(myRootCtx, "player", c)
//	go processProtocol(subCtx, wg, chRead, chWrite)
//	// 收数据
//	wg.Add(1)
//	go func() {
//		defer wg.Done()
//		// 读消息
//		for {
//			select {
//			case <-myRootCtx.Done():
//				log.Printf("conn %d eixt\n", index)
//				return
//			default:
//				err := wsReadBuff(conn, chRead)
//				if err != nil {
//					log.Printf("Read protocol failed, index: %v, err: %v\n", index, err)
//					myCancel()
//				}
//			}
//		}
//	}()
//	// 发数据
//	for {
//		select {
//		case <-myRootCtx.Done():
//			conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
//			log.Printf("conn %d eixt\n", index)
//			return
//		case pbBuff := <-chWrite:
//			if len(pbBuff) >= protocols.HeadSize {
//				//var buff bytes.Buffer
//				enMessage := base64.StdEncoding.EncodeToString(pbBuff)
//				//fmt.Println(enMessage)
//				enMessage += postfix
//				err = conn.WriteMessage(websocket.TextMessage, []byte(enMessage))
//				if err != nil {
//					log.Println("write failed, err:", err)
//					return
//				}
//			}else {
//				myCancel()
//			}
//		}
//	}
//}
//
//func wsReadBuff(conn *websocket.Conn, chRead chan<- *protocols.Protocol) error {
//	_, message, err := conn.ReadMessage()
//	if err != nil {
//		log.Println("read failed, err:", err)
//		return err
//	}
//	newMessage := strings.Trim(string(message), postfix)
//	// 解码
//	deBuff, err:= base64.StdEncoding.DecodeString(newMessage)
//	if err != nil {
//		log.Println("base64 decoding failed, err:", err)
//		return err
//	}
//	br := bytes.NewReader(deBuff)
//	for br.Len() > 0 {
//		ptData := new(protocols.Protocol)
//		err = binary.Read(br, binary.LittleEndian, &ptData.Head)
//		if err != nil {
//			log.Println("binary.read failed, error is ", err)
//			return err
//		}
//		if ptData.Head.Len > protocols.HeadSize {
//			leftBuff := make([]byte, ptData.Head.Len-protocols.HeadSize)
//			if _, err := io.ReadFull(br, leftBuff); err != nil {
//				fmt.Println("read protocol content failed, error is ", err)
//				return err
//			}
//			ptData.Content.Write(leftBuff)
//		}
//		chRead <- ptData
//	}
//	return nil
//}

