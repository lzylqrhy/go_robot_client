package protocols

import "github/go-robot/util"

const (
	EnterRoomCode = 0x301	// 进入房间
	SceneInfoCode = 0x305	// 请求场景信息
	PlayerSeatCode = 0x306	// 鱼池座位列表信息
	TransmitCode = 0x318	// 转发行为
)

type C2SEnterRoom struct {
	RoomID uint32
	ChannelID uint32
}

func (p *C2SEnterRoom) Bytes() []byte {
	var pb Protocol
	pb.SetCmd(EnterRoomCode)
	pb.AppendNumber(p.RoomID)
	pb.AppendNumber(p.ChannelID)
	return pb.Bytes()
}

type S2CEnterRoom struct {
	RoomID uint32
	Result uint8
}

func (p *S2CEnterRoom) Parse(pb *Protocol) {
	err := pb.GetNumber(p)
	util.CheckError(err)
}

type C2SGetSceneInfo struct {

}

func (p *C2SGetSceneInfo) Bytes() []byte {
	var pb Protocol
	pb.SetCmd(SceneInfoCode)
	return pb.Bytes()
}

type S2CGetSceneInfo struct {
	BGID uint32
	SevTime uint64
	Buffs []struct {
		BuffType uint32
		Duration uint32
	}
}

func (p *S2CGetSceneInfo) Parse(pb *Protocol) {
	err := pb.GetNumber(&p.BGID)
	util.CheckError(err)
	err = pb.GetNumber(&p.SevTime)
	util.CheckError(err)
	var count uint8
	err = pb.GetNumber(&count)
	util.CheckError(err)
	//p.Buffs
	for i:=uint8(0); i < count; i++{
		err = pb.GetNumber(&p.Buffs[i])
		util.CheckError(err)
	}
}

type C2STransmitActivity struct {
	Activity string
}

func (p *C2STransmitActivity) Bytes() []byte {
	var pb Protocol
	pb.SetCmd(TransmitCode)
	pb.AppendStringUint8(p.Activity)
	return pb.Bytes()
}