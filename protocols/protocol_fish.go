package protocols

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func PacketPing() []byte {
	header := ProtocolHeader{3, 4}
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.LittleEndian, &header)
	if err != nil {
		fmt.Println("packet protocol failed, header is ", header)
	}
	return buff.Bytes()
}
