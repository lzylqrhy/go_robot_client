package net

import (
	"bytes"
	"context"
	"encoding/binary"
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

func NewWSConnect(sAddr string) MyDialer {
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
	log.Printf("web socket dial %s successfully", d.sAddr)
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
			log.Println("send on disconnected message")
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
				log.Println("conn exit")
				return
			case data := <-d.chWrite:
				if data == nil {
					return
				}
				/* string传输方式
				enMessage := base64.StdEncoding.EncodeToString(data)
				enMessage += postfix
				err := d.conn.WriteMessage(websocket.TextMessage, []byte(enMessage))
				 */
				err := d.conn.WriteMessage(websocket.BinaryMessage, data)
				if err != nil {
					log.Println("write failed, err:", err)
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
	headBuff := make([]byte, protocols.HeadSize)
	leftCount := uint16(0)
	for {
		_, message, err := d.conn.ReadMessage()
		if err != nil {
			log.Println("read failed, err:", err)
			break
		}
		/* string传输方式
		// 去掉后缀
		newMessage := strings.Trim(string(message), postfix)
		// 解码
		deBuff, err:= base64.StdEncoding.DecodeString(newMessage)
		if err != nil {
			log.Println("base64 decoding failed, err:", err)
			d.Disconnect()
			continue
		}
		br := bytes.NewReader(deBuff)
		 */
		br := bytes.NewReader(message)
		for {
			if nil == ptData {
				if 0 == br.Len() {
					break
				}
				if br.Len() + int(leftCount) < protocols.HeadSize {
					n, err := br.Read(headBuff[leftCount:])
					if err != nil {
						log.Println("read head buff failed, error is ", err)
						return
					}
					headBuff = headBuff[:leftCount+uint16(n)]
					leftCount = protocols.HeadSize - (leftCount+uint16(n))
					break
				}
				ptData = new(protocols.Protocol)
				err = binary.Read(br, binary.LittleEndian, &ptData.Head)
				if err != nil {
					log.Println("binary.read failed, error is ", err)
					return
				}
				if ptData.Head.Len > protocols.HeadSize {
					if 0 == br.Len() {
						leftCount = ptData.Head.Len - protocols.HeadSize
						break
					}
					leftBuff := make([]byte, ptData.Head.Len-protocols.HeadSize)
					if n, err := io.ReadFull(br, leftBuff); err != nil {
						if err == io.ErrUnexpectedEOF {
							leftCount = ptData.Head.Len - protocols.HeadSize - uint16(n)
							ptData.Content.Write(leftBuff)
							break
						}else {
							log.Println("read protocol content failed, error is ", err)
							return
						}
					}
					ptData.Content.Write(leftBuff)
				}
				d.chRead <- ptData
				ptData = nil
				leftCount = 0
			}else {
				if br.Len() < int(leftCount) {
					leftBuff := make([]byte, br.Len())
					n, err := br.Read(leftBuff)
					if err != nil {
						log.Println("read head buff failed, error is ", err)
						return
					}
					leftCount -= uint16(n)
					break
				}
				leftBuff := make([]byte, leftCount)
				if _, err := io.ReadFull(br, leftBuff); err != nil {
					log.Println("read protocol content failed, error is ", err)
				}
				ptData.Content.Write(leftBuff)
				d.chRead <- ptData
				ptData = nil
				leftCount = 0
			}
		}
	}
}

