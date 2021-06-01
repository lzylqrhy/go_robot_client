package protocols

import (
	"github/go-robot/util"
)

const (
	//FruitEnterRoomCode       = 0x10 // 进入房间
	//FruitPlayerCode          = 0x501 // 玩家个人信息
	FruitJoinRoomCode		= 0x590	// 加入房间
	FruitPlayerCode			= 0x596	// 玩家个人信息
	FruitPlayCode			= 0x531	// 开始游戏、游戏结果
)

type FruitS2CPlayerInfo struct {
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

func (p *FruitS2CPlayerInfo) Parse(pb *Protocol) {
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
type FruitC2SJoinRoom struct {
	GameID uint8
}

func (p *FruitC2SJoinRoom) Bytes() []byte {
	var pb Protocol
	pb.SetCmd(FruitJoinRoomCode)
	pb.AppendNumber(p.GameID)
	return pb.Bytes()
}

type FruitS2CJoinRoom struct {
	Result uint8
}

func (p *FruitS2CJoinRoom) Parse(pb *Protocol) {
	err := pb.GetNumber(p)
	util.CheckError(err)
}

// 开始游戏
type FruitC2SPlay struct {
	Line uint8
	Amount uint32
	GameID uint8
}

func (p *FruitC2SPlay) Bytes() []byte {
	var pb Protocol
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

func (line *fruitLine) Parse(pb *Protocol) {
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

func (re *fruitResult) Parse(pb *Protocol) {
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

type FruitS2CPlayResult struct {
	Result uint8
	ReInfo []fruitResult
}

func (p *FruitS2CPlayResult) Parse(pb *Protocol) {
	var line int32
	util.CheckError(pb.GetNumber(&line))
	util.CheckError(pb.GetNumber(&p.Result))
	var count uint8
	util.CheckError(pb.GetNumber(&count))
	p.ReInfo = make([]fruitResult, count)
	for i := uint8(0); i < count; i++ {
		p.ReInfo[i].Parse(pb)
	}
}
