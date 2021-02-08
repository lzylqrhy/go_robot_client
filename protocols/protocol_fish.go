package protocols

import "github/go-robot/util"

const (
	EnterRoomCode = 0x301	// 进入房间
	SceneInfoCode = 0x305	// 请求场景信息
	PlayerSeatCode = 0x306	// 鱼池座位列表信息
	FishListCode = 0x307	// 鱼池鱼列表信息
	BulletListCode = 0x308	// 鱼池子弹列表信息
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

type PtBuff struct {
	BuffType uint32
	Data 	 uint32
	Duration uint32
}

type S2CGetSceneInfo struct {
	BGImgID uint32
	ServerTime uint64
	Buffs []PtBuff
}

func (p *S2CGetSceneInfo) Parse(pb *Protocol) {
	util.CheckError(pb.GetNumber(&p.BGImgID))
	util.CheckError(pb.GetNumber(&p.ServerTime))
	var count uint8
	util.CheckError(pb.GetNumber(&count))
	p.Buffs = make([]PtBuff, count)
	for i:=uint8(0); i < count; i++{
		util.CheckError(pb.GetNumber(&p.Buffs[i].BuffType))
		util.CheckError(pb.GetNumber(&p.Buffs[i].Duration))
	}
}

type PtSeat struct {
	CharID    uint32
	Figure    string
	Nick      string
	Currency  uint64
	LV        uint8
	SeatID    uint8
	CannonID  uint32
	CaliberLV uint8
	Caliber   uint32
	Status    uint16
	Buffs     []PtBuff
}

type S2CSeatsInfo struct {
	Players []PtSeat
}

func (p *S2CSeatsInfo) Parse(pb *Protocol) {
	var seatCount uint8
	util.CheckError(pb.GetNumber(&seatCount))
	p.Players = make([]PtSeat, seatCount)
	var err error
	var buffCount uint8
	for i := uint8(0); i < seatCount; i++ {
		player := &p.Players[i]
		util.CheckError(pb.GetNumber(&player.CharID))
		player.Figure, err = pb.GetStringUint8()
		util.CheckError(err)
		player.Nick, err = pb.GetStringUint8()
		util.CheckError(err)
		util.CheckError(pb.GetNumber(&player.Currency))
		util.CheckError(pb.GetNumber(&player.LV))
		util.CheckError(pb.GetNumber(&player.SeatID))
		util.CheckError(pb.GetNumber(&player.CannonID))
		util.CheckError(pb.GetNumber(&player.CaliberLV))
		util.CheckError(pb.GetNumber(&player.Caliber))
		util.CheckError(pb.GetNumber(&player.Status))
		util.CheckError(pb.GetNumber(&buffCount))
		player.Buffs = make([]PtBuff, buffCount)
		for j := uint8(0); j < buffCount; j++ {
			util.CheckError(pb.GetNumber(&player.Buffs[i].BuffType))
			util.CheckError(pb.GetNumber(&player.Buffs[i].Duration))
		}
	}
}

type PtFish struct {
	Serial   uint32
	FishID   uint32
	PathID   uint32
	Speed    uint32
	OffsetX  uint32
	OffsetY  uint32
	OffsetZ  uint32
	BornTime uint64
}

type S2CFishList struct {
	FishList []PtFish
}

func (p *S2CFishList) Parse(pb *Protocol) {
	var count uint16
	util.CheckError(pb.GetNumber(&count))
	p.FishList = make([]PtFish, count)
	for i := uint16(0); i < count; i++ {
		util.CheckError(pb.GetNumber(&p.FishList[i]))
	}
}

type PtBullet struct {
	Serial   uint32
	OriginID uint32
	SeatID   uint32
	CharID   uint32
	SkinID   uint32
	Radian   float32
	BornTime uint64
	Buffs    []PtBuff
}

type S2CBulletList struct {
	BulletList []PtBullet
}

func (p *S2CBulletList) Parse(pb *Protocol) {
	var count uint8
	util.CheckError(pb.GetNumber(&count))
	p.BulletList = make([]PtBullet, count)
	var buffCount uint8
	for i := uint8(0); i < count; i++ {
		bullet := &p.BulletList[i]
		util.CheckError(pb.GetNumber(&bullet.Serial))
		util.CheckError(pb.GetNumber(&bullet.OriginID))
		util.CheckError(pb.GetNumber(&bullet.SeatID))
		util.CheckError(pb.GetNumber(&bullet.CharID))
		util.CheckError(pb.GetNumber(&bullet.SkinID))
		util.CheckError(pb.GetNumber(&bullet.Radian))
		util.CheckError(pb.GetNumber(&bullet.BornTime))
		util.CheckError(pb.GetNumber(&buffCount))
		bullet.Buffs = make([]PtBuff, buffCount)
		for j := uint8(0); j < buffCount; j++ {
			util.CheckError(pb.GetNumber(&bullet.Buffs[i].BuffType))
			util.CheckError(pb.GetNumber(&bullet.Buffs[i].Data))
		}
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