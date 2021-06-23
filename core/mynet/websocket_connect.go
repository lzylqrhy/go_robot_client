package mynet

import (
	"bytes"
	"context"
	"encoding/binary"
	"github.com/gorilla/websocket"
	"github/go-robot/core/protocol"
	"io"
	"log"
	"net/url"
	"sync"
)

//const postfix = "\r\n\r\n"

type WSDialer struct {
	conn *websocket.Conn
	chRead chan *protocol.Protocol
	chWrite chan []byte
	sAddr string
	ctx context.Context
}

func NewWSConnect(sAddr string) MyDialer {
	d := new(WSDialer)
	u := url.URL{Scheme: "ws", Host: sAddr}
	d.sAddr = u.String()
	d.chRead = make(chan *protocol.Protocol, 100)
	d.chWrite = make(chan []byte, 100)
	return d
}

func (d *WSDialer) connect() bool {
	if d.conn != nil {
		return true
	}
	var err error
	d.conn, _, err = websocket.DefaultDialer.Dial(d.sAddr, nil)
	if err != nil {
		log.Println("web socket dial failed, err:", err)
		return false
	}
	log.Printf("web socket dial %s successfully", d.sAddr)
	return true
}

func (d *WSDialer) Disconnect() {
	if d.conn != nil {
		// 单0表示关闭
		d.SendPacket([]byte{0})
		Break:
		for {
			select {
			case pb := <-d.chRead:
				if 0 == pb.Head.Cmd {
					close(d.chWrite)
					close(d.chRead)
					break Break
				}
			}
		}
		log.Println("web socket disconnect successfully")
	}
}

func (d *WSDialer) close() {
	if d.conn != nil {
		// 断开连接
		log.Println("close web socket connect")
		err := d.conn.Close()
		if err != nil {
			log.Println("web socket close socket failed, err:", err)
		}
		d.conn = nil
	}
}

func (d *WSDialer) SendPacket(data []byte) bool {
	if d.conn == nil {
		log.Println("send failed, the dialer is offline")
		return false
	}
	d.chWrite<- data
	return true
}

func (d *WSDialer) ReadPacket() <-chan *protocol.Protocol {
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
		var pd protocol.Protocol
		pd.Head.Cmd = 0
		pd.Head.Len = 0
		d.chRead <- &pd
		// 读消息
		go func() {
			defer func() {
				d.close()
				// 连接断开
				var pd protocol.Protocol
				pd.Head.Cmd = 0
				pd.Head.Len = 1
				d.chRead <- &pd
				// 关闭发数据的loop，一定要使用
				d.chWrite <- nil
			}()
			// 读消息
			d.read()
		}()
		isDisconnect := false
		// 发数据
		for {
			select {
			case <-ctx.Done():
				if !isDisconnect {
					go d.Disconnect()
					isDisconnect = true
				}
			case data := <-d.chWrite:
				if data == nil {
					return
				}
				if d.conn == nil {
					break
				}
				if bytes.Equal(data, []byte{0}) {
					err := d.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
					if err != nil {
						log.Println("web socket soft close failed, err:", err)
					}
					break
				}
				/* string传输方式
				enMessage := base64.StdEncoding.EncodeToString(data)
				enMessage += postfix
				err := d.conn.WriteMessage(websocket.TextMessage, []byte(enMessage))
				 */
				//log.Println("write msg:", data[:4], "\n", string(debug.Stack()))
				err := d.conn.WriteMessage(websocket.BinaryMessage, data)
				if err != nil {
					log.Println("write failed, err:", err)
					return
				}
				//log.Println("write msg ok :", data[:4])
			}
		}
	}()
	return true
}

func (d *WSDialer) read() {
	// 读消息
	var ptData *protocol.Protocol
	headBuff := make([]byte, protocol.HeadSize)
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
				if br.Len() + int(leftCount) < protocol.HeadSize {
					n, err := br.Read(headBuff[leftCount:])
					if err != nil {
						log.Println("read head buff failed, error is ", err)
						return
					}
					headBuff = headBuff[:leftCount+uint16(n)]
					leftCount = protocol.HeadSize - (leftCount+uint16(n))
					break
				}
				ptData = new(protocol.Protocol)
				err = binary.Read(br, binary.LittleEndian, &ptData.Head)
				if err != nil {
					log.Println("binary.read failed, error is ", err)
					return
				}
				if ptData.Head.Len > protocol.HeadSize {
					if 0 == br.Len() {
						leftCount = ptData.Head.Len - protocol.HeadSize
						break
					}
					leftBuff := make([]byte, ptData.Head.Len-protocol.HeadSize)
					if n, err := io.ReadFull(br, leftBuff); err != nil {
						if err == io.ErrUnexpectedEOF {
							leftCount = ptData.Head.Len - protocol.HeadSize - uint16(n)
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

