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

type Protocol struct {
	Header  ProtocolHeader
	Content bytes.Buffer
}

func (pb *Protocol)Bytes() []byte {
	buff := new(bytes.Buffer)
	if err := binary.Write(buff, binary.LittleEndian, pb.Header); err != nil {
		fmt.Println("packet protocol failed, header is ", pb.Header)
		return nil
	}
	buff.Write(pb.Content.Bytes())
	return buff.Bytes()
}

func (pb *Protocol)AppendNumber(value interface{}) {
	if err := binary.Write(&pb.Content, binary.LittleEndian, value); err != nil {
		fmt.Println("append number failed, err: ", err)
	}
}

func (pb *Protocol)AppendString(lenByteNum uint8, value string) {
	if err := binary.Write(&pb.Content, binary.LittleEndian, lenByteNum); err != nil {
		fmt.Println("append number failed, err: ", err)
	}
	pb.Content.WriteString(value)
}