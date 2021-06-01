package protocols

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"math"
	"reflect"
)

const HeadSize = 4	// 协议头长度

// 协议头
type ProtocolHeader struct {
	Cmd uint16
	Len uint16
}

// 协议数据类
type Protocol struct {
	Head ProtocolHeader
	Content bytes.Buffer
}

// 设置协议号
func (p *Protocol)SetCmd(cmd uint16) {
	p.Head.Cmd = cmd
}

// 将协议序列化为字节切片
func (p *Protocol)Bytes() []byte {
	p.Head.Len = uint16(p.Content.Len() + HeadSize)
	buff := new(bytes.Buffer)
	if err := binary.Write(buff, binary.LittleEndian, &p.Head); err != nil {
		log.Println("append number failed, err: ", err)
	}
	buff.Write(p.Content.Bytes())
	return buff.Bytes()
}

func (p *Protocol)AppendNumber(value interface{}) {
	if err := binary.Write(&p.Content, binary.LittleEndian, value); err != nil {
		log.Println("append number failed, err: ", err)
	}
}

func (p *Protocol)AppendStringUint8(value string) {
	length := len(value)
	if length >= math.MaxUint8 {
		panic("length of string is greater than max uint8")
	}
	p.appendString(uint8(length), value)
}

func (p *Protocol)AppendStringUint16(value string) {
	length := len(value)
	if length >= math.MaxUint16 {
		panic("length of string is greater than max uint16")
	}
	p.appendString(uint16(length), value)
}

func (p *Protocol)appendString(length interface{}, value string) {
	p.AppendNumber(length)
	p.Content.WriteString(value)
}

func (p *Protocol)GetNumber(value interface{}) error {
	if reflect.TypeOf(value).Kind() != reflect.Ptr {
		log.Fatalln("value must be pointer type")
	}
	return binary.Read(&p.Content, binary.LittleEndian, value)
}

func (p *Protocol)GetStringUint8() (string, error) {
	var length uint8
	return p.getString(&length)
}

func (p *Protocol)GetStringUint16() (string, error) {
	var length uint16
	return p.getString(&length)
}

func (p *Protocol)getString(length interface{}) (string, error) {
	if err := p.GetNumber(length); err != nil {
		return "", err
	}
	var n int
	switch sz := length.(type) {
	case *uint8:
		n = int(*sz)
	case *uint16:
		n = int(*sz)
	default:
		return "", errors.New("the type of length of string is invalid")
	}
	buff := p.Content.Next(n)
	if len(buff) < n {
		return "", errors.New("buff can't fill the string")
	}
	return string(buff), nil
}

