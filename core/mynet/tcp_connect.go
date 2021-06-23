package mynet

import (
	"bytes"
	"context"
	"encoding/binary"
	"github/go-robot/core/protocol"
	"io"
	"log"
	"net"
	"sync"
)

type TCPDialer struct {
	conn net.Conn
	chRead chan *protocol.Protocol
	chWrite chan []byte
	sAddr string
	ctx context.Context
}

func NewTCPConnect(sAddr string) MyDialer {
	d := new(TCPDialer)
	d.sAddr = sAddr
	d.chRead = make(chan *protocol.Protocol, 100)
	d.chWrite = make(chan []byte, 100)
	return d
}

func (d *TCPDialer) connect() bool {
	if d.conn != nil {
		return true
	}
	var err error
	d.conn, err = net.Dial(TCP, d.sAddr)
	if err != nil {
		log.Println("tcp socket dial failed, err:", err)
		return false
	}
	log.Printf("tcp socket dial %s successfully", d.sAddr)
	return true
}

func (d *TCPDialer) Disconnect() {
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
		log.Println("tcp socket disconnect successfully")
	}
}

func (d *TCPDialer) close() {
	if d.conn != nil {
		// 断开连接
		log.Println("close tcp socket connect")
		err := d.conn.Close()
		if err != nil {
			log.Println("tcp socket close socket failed, err:", err)
		}
		d.conn = nil
	}
}

func (d *TCPDialer) SendPacket(data []byte) bool {
	if d.conn == nil {
		log.Println("send failed, the dialer is offline")
		return false
	}
	d.chWrite<- data
	return true
}

func (d *TCPDialer) ReadPacket() <-chan *protocol.Protocol {
	return d.chRead
}

func (d *TCPDialer) Run(ctx context.Context, wg *sync.WaitGroup) bool {
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
					d.conn.Close()
					break
				}
				//log.Println("write msg:", data[:4], "\n", string(debug.Stack()))
				_, err := d.conn.Write(data)
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

func (d *TCPDialer) read() {
	headBuff := make([]byte, protocol.HeadSize)
	for {
		if _, err := io.ReadFull(d.conn, headBuff); err != nil {
			log.Println("read protocol header failed, error is ", err)
			break
		}
		ptData := new(protocol.Protocol)
		r := bytes.NewReader(headBuff)
		err := binary.Read(r, binary.LittleEndian, &ptData.Head)
		if err != nil {
			log.Println("binary.read failed, error is ", err)
			break
		}
		if ptData.Head.Len > protocol.HeadSize {
			leftBuff := make([]byte, ptData.Head.Len-protocol.HeadSize)
			if _, err := io.ReadFull(d.conn, leftBuff); err != nil {
				log.Println("read protocol content failed, error is ", err)
				break
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
		d.chRead <- ptData
	}
}
