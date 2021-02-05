package protocols

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

type ProtocolHeader struct {
	Cmd uint16
	Len uint16
}

type Protocol struct {
	Content bytes.Buffer
}

func (pb *Protocol)Bytes() []byte {
	return pb.Content.Bytes()
}

func (pb *Protocol)AppendNumber(value interface{}) {
	if err := binary.Write(&pb.Content, binary.LittleEndian, value); err != nil {
		fmt.Println("append number failed, err: ", err)
	}
}

func (pb *Protocol)AppendStringUint8(value string) {
	length := len(value)
	if length >= math.MaxUint8 {
		panic("length of string is greater than max uint8")
	}
	pb.appendString(uint8(length), value)
}

func (pb *Protocol)AppendStringUint16(value string) {
	length := len(value)
	if length >= math.MaxUint16 {
		panic("length of string is greater than max uint16")
	}
	pb.appendString(uint16(length), value)
}

func (pb *Protocol)appendString(length interface{}, value string) {
	pb.AppendNumber(length)
	pb.Content.WriteString(value)
}

func (pb *Protocol)GetNumber(ref interface{}) error {
	return binary.Read(&pb.Content, binary.LittleEndian, ref)
}

func (pb *Protocol)GetStringUint8() (string, error) {
	var length uint8
	return pb.getString(&length)
}

func (pb *Protocol)GetStringUint16() (string, error) {
	var length uint16
	return pb.getString(&length)
}

func (pb *Protocol)getString(length interface{}) (string, error) {
	if err := pb.GetNumber(length); err != nil {
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
	buff := pb.Content.Next(n)
	if len(buff) < n {
		return "", errors.New("buff can't fill the string")
	}
	return string(buff), nil
}