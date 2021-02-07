package protocols

import "github/go-robot/util"

const (
	PingCode = 0x1
	SyncTimeCode = 0x3
	C2SLoginCode = 0xC2
	S2CLoginCode = 0x10
	ResLoadedCode = 0xB6
)

type C2SPing struct {
	TimeStamp uint32
}

func (p *C2SPing) Bytes() []byte {
	var pb Protocol
	pb.SetCmd(PingCode)
	return pb.Bytes()
}

type C2SSyncTime struct {

}

func (p *C2SSyncTime) Bytes() []byte {
	var pb Protocol
	pb.SetCmd(SyncTimeCode)
	return pb.Bytes()
}

type S2CSyncTime struct {
	Year uint16
	Month uint8
	Day uint8
	Hour uint8
	Minute uint8
	Second uint8
	TimeStamp uint32
}

func (p *S2CSyncTime) Parse(pb *Protocol) {
	err := pb.GetNumber(p)
	util.CheckError(err)
}

type C2SLogin struct {
	IsChildGame bool
	Token string
}

func (p *C2SLogin) Bytes() []byte {
	var pb Protocol
	pb.SetCmd(C2SLoginCode)
	pb.AppendNumber(p.IsChildGame)
	pb.AppendStringUint8(p.Token)
	return pb.Bytes()
}

type S2CLogin struct {
	Status uint8
	TimeStamp uint32
}

func (p *S2CLogin) Parse(pb *Protocol) {
	err := pb.GetNumber(p)
	util.CheckError(err)
}

type C2SResourceLoaded struct {

}

func (p *C2SResourceLoaded) Bytes() []byte {
	var pb Protocol
	pb.SetCmd(ResLoadedCode)
	return pb.Bytes()
}