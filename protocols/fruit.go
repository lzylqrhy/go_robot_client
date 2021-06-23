package protocols

import (
	"github/go-robot/core/protocol"
	"github/go-robot/util"
)

const (
	FruitJoinRoomCode		= 0x590	// 加入房间
	FruitPlayerCode			= 0x596	// 玩家个人信息
	FruitPlayCode			= 0x531	// 开始游戏、游戏结果
)

type S2CFruitPlayerInfo struct {
	CharID     uint32
	PID        uint32
	RegTime    uint32
	Gender     uint8
	A          uint32
	LV         uint8
	B          uint8
	C          uint32
	VipLV      uint32
	Exp        uint32
	UpgradeExp uint32
	Nick       string
	Figure     string
	addr       string
}

func (p *S2CFruitPlayerInfo) Parse(pb *protocol.Protocol) {
	var err error
	util.CheckError(pb.GetNumber(&p.CharID))
	util.CheckError(pb.GetNumber(&p.PID))
	util.CheckError(pb.GetNumber(&p.RegTime))
	util.CheckError(pb.GetNumber(&p.Gender))
	util.CheckError(pb.GetNumber(&p.A))
	util.CheckError(pb.GetNumber(&p.LV))
	util.CheckError(pb.GetNumber(&p.B))
	util.CheckError(pb.GetNumber(&p.C))
	util.CheckError(pb.GetNumber(&p.VipLV))
	util.CheckError(pb.GetNumber(&p.Exp))
	util.CheckError(pb.GetNumber(&p.UpgradeExp))
	p.Nick, err = pb.GetStringUint8()
	util.CheckError(err)
	p.Figure, err = pb.GetStringUint16()
	util.CheckError(err)
	p.addr, err = pb.GetStringUint8()
	util.CheckError(err)
}

// 加入房间
type C2SFruitJoinRoom struct {
	GameID uint8
}

func (p *C2SFruitJoinRoom) Bytes() []byte {
	var pb protocol.Protocol
	pb.SetCmd(FruitJoinRoomCode)
	pb.AppendNumber(p.GameID)
	return pb.Bytes()
}

type S2CFruitJoinRoom struct {
	Result uint8
}

func (p *S2CFruitJoinRoom) Parse(pb *protocol.Protocol) {
	err := pb.GetNumber(p)
	util.CheckError(err)
}

// 开始游戏
type C2SFruitPlay struct {
	Line uint8
	Amount uint32
	GameID uint8
}

func (p *C2SFruitPlay) Bytes() []byte {
	var pb protocol.Protocol
	pb.SetCmd(FruitPlayCode)
	pb.AppendNumber(p)
	return pb.Bytes()
}

type fruitIcon struct {
	Index, Icon uint8
}

type fruitLine struct {
	ID uint8
	Point []uint8
	Multiple uint32
	IsFree uint8
}

func (line *fruitLine) Parse(pb *protocol.Protocol) {
	util.CheckError(pb.GetNumber(&line.ID))
	// 连续点的数量
	var lineNum uint8
	util.CheckError(pb.GetNumber(&lineNum))
	line.Point = make([]uint8, lineNum)
	for j := uint8(0); j < lineNum; j++ {
		util.CheckError(pb.GetNumber(&line.Point[j]))
	}
	util.CheckError(pb.GetNumber(&line.Multiple))
	util.CheckError(pb.GetNumber(&line.IsFree))
}

type fruitResult struct {
	Icon []fruitIcon
	Lines []fruitLine
	Amount uint32
}

func (re *fruitResult) Parse(pb *protocol.Protocol) {
	var num uint8
	util.CheckError(pb.GetNumber(&num))
	re.Icon = make([]fruitIcon, num)
	for i := uint8(0); i < num; i++ {
		util.CheckError(pb.GetNumber(&re.Icon[i]))
	}
	util.CheckError(pb.GetNumber(&num))
	re.Lines = make([]fruitLine, num)
	for i := uint8(0); i < num; i++ {
		re.Lines[i].Parse(pb)
	}
	util.CheckError(pb.GetNumber(&re.Amount))
}

type S2CFruitPlayResult struct {
	Line uint32
	Result uint8
	ReInfo []fruitResult
}

func (p *S2CFruitPlayResult) Parse(pb *protocol.Protocol) {
	util.CheckError(pb.GetNumber(&p.Line))
	util.CheckError(pb.GetNumber(&p.Result))
	var count uint8
	util.CheckError(pb.GetNumber(&count))
	p.ReInfo = make([]fruitResult, count)
	for i := uint8(0); i < count; i++ {
		p.ReInfo[i].Parse(pb)
	}
}
