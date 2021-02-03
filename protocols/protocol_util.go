package protocols

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type ProtocolHeader struct {
	Cmd uint16
	Len uint16
}

func (header *ProtocolHeader)Bytes() *bytes.Buffer {
	buff := new(bytes.Buffer)
	if err := binary.Write(buff, binary.LittleEndian, header); err != nil {
		fmt.Println("packet protocol failed, header is ", header)
		return buff
	}
	return buff
}

type Protocol struct {
	Header  ProtocolHeader
	Content bytes.Buffer
}

func (pb *Protocol)Bytes() []byte {
	pb.Header.Bytes().Write(pb.Content.Bytes())
	return pb.Bytes()
}

