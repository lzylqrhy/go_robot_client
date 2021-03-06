package protocols

import (
	"encoding/json"
	"github/go-robot/core/protocol"
	"github/go-robot/util"
)

const (
	PlayerCode          = 0x101  // 玩家个人信息
	EnterHallOrRoomCode = 0x103  // 登录成功后进入大厅或房间
	ReadPacketInfoCode  = 0x105  // 红包抽奖信息
	DrawReadPacketCode  = 0x106  // 开红包
	RoomListCode        = 0x300  // 房间列表
	FishEnterRoomCode   = 0x301  // 进入房间
	SceneInfoCode       = 0x305  // 请求场景信息
	PlayerSeatCode      = 0x306  // 鱼池座位列表信息
	FishListCode        = 0x307  // 鱼池鱼列表信息
	BulletListCode      = 0x308  // 鱼池子弹列表信息
	FireCode            = 0x309  // 开火
	SwitchCaliberCode   = 0x30A  // 开火
	HitFishCode         = 0x30B  // 命中鱼
	SyncFishBoom        = 0x30D  // 同步鱼潮
	GenerateFish        = 0x30E  // 生成鱼
	TransmitCode        = 0x318  // 转发行为
	LaunchMissileCode		= 0x319	 // 发射导弹
	PoseidonStatusCode  = 0x010C // 波塞冬状态
	HitPoseidonCode    = 0x010B  // 命中波塞冬
)

type S2CPlayerInfo struct {
	CharID     uint32
	Figure     string
	Nick       string
	LV         uint8
	Exp        uint32
	UpgradeExp uint32
}

func (p *S2CPlayerInfo) Parse(pb *protocol.Protocol) {
	var err error
	util.CheckError(pb.GetNumber(&p.CharID))
	p.Figure, err = pb.GetStringUint8()
	util.CheckError(err)
	p.Nick, err = pb.GetStringUint8()
	util.CheckError(err)
	util.CheckError(pb.GetNumber(&p.LV))
	util.CheckError(pb.GetNumber(&p.Exp))
	util.CheckError(pb.GetNumber(&p.UpgradeExp))
}

type S2CEnterHallOrRoom struct {
	RoomID uint32
}

func (p *S2CEnterHallOrRoom) Parse(pb *protocol.Protocol) {
	err := pb.GetNumber(p)
	util.CheckError(err)
}

type Room struct {
	RoomID uint32
	Type uint8
	Status uint8
	CostType uint32
	MinScore uint32
	MaxScore uint32
	MinCannon uint32
	MaxCannon uint32
}

type S2CRoomList struct {
	Rooms []Room
}

func (p *S2CRoomList) Parse(pb *protocol.Protocol) {
	var count uint8
	util.CheckError(pb.GetNumber(&count))
	p.Rooms = make([]Room, count)
	for i := 0; i < int(count); i++ {
		room := &p.Rooms[i]
		util.CheckError(pb.GetNumber(&room.RoomID))
		util.CheckError(pb.GetNumber(&room.Type))
		util.CheckError(pb.GetNumber(&room.Status))
		util.CheckError(pb.GetNumber(&room.CostType))
		util.CheckError(pb.GetNumber(&room.MinScore))
		util.CheckError(pb.GetNumber(&room.MaxScore))
		util.CheckError(pb.GetNumber(&room.MinCannon))
		util.CheckError(pb.GetNumber(&room.MaxCannon))
	}
}

type C2SFishEnterRoom struct {
	RoomID uint32
	ChannelID uint32
}

func (p *C2SFishEnterRoom) Bytes() []byte {
	var pb protocol.Protocol
	pb.SetCmd(FishEnterRoomCode)
	pb.AppendNumber(p.RoomID)
	pb.AppendNumber(p.ChannelID)
	return pb.Bytes()
}

type S2CFishEnterRoom struct {
	RoomID uint32
	Result uint8
}

func (p *S2CFishEnterRoom) Parse(pb *protocol.Protocol) {
	err := pb.GetNumber(p)
	util.CheckError(err)
}

type C2SGetSceneInfo struct {

}

func (p *C2SGetSceneInfo) Bytes() []byte {
	var pb protocol.Protocol
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
	ServerTime float64
	Buffs []PtBuff
}

func (p *S2CGetSceneInfo) Parse(pb *protocol.Protocol) {
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
	VIP		  uint8
	Currency  float64
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

func (p *S2CSeatsInfo) Parse(pb *protocol.Protocol) {
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
		util.CheckError(pb.GetNumber(&player.VIP))
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
			util.CheckError(pb.GetNumber(&player.Buffs[j].BuffType))
			util.CheckError(pb.GetNumber(&player.Buffs[j].Duration))
		}
	}
}

type PtFish struct {
	Serial   uint32
	KindID   uint32
	PathID   uint32
	Speed    uint32
	OffsetX  uint32
	OffsetY  uint32
	OffsetZ  uint32
	BornTime float64
	SwamTime uint32
}

type S2CFishList struct {
	FishList []PtFish
}

func (p *S2CFishList) Parse(pb *protocol.Protocol) {
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
	SeatID   uint8
	CharID   uint32
	SkinID   uint32
	Radian   float32
	BornTime float64
	Buffs    []PtBuff
}

type S2CBulletList struct {
	BulletList []PtBullet
}

func (p *S2CBulletList) Parse(pb *protocol.Protocol) {
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
			util.CheckError(pb.GetNumber(&bullet.Buffs[j].BuffType))
			util.CheckError(pb.GetNumber(&bullet.Buffs[j].Data))
		}
	}
}

type C2SFire struct {
	OriginID uint32
	Radian float32
	TargetSerial uint32
}

func (p *C2SFire) Bytes() []byte {
	var pb protocol.Protocol
	pb.SetCmd(FireCode)
	pb.AppendNumber(p)
	return pb.Bytes()
}

type S2CFire struct {
	Result   uint8
	Serial   uint32
	OriginID uint32
	Cost     uint32
	Currency float64
	SeatID   uint8
	CharID   uint32
	SkinID   uint32
	Radian   float32
	BornTime float64
	Buffs    []PtBuff
}

func (p *S2CFire) Parse(pb *protocol.Protocol) {
	util.CheckError(pb.GetNumber(&p.Result))
	util.CheckError(pb.GetNumber(&p.Serial))
	util.CheckError(pb.GetNumber(&p.OriginID))
	util.CheckError(pb.GetNumber(&p.Cost))
	util.CheckError(pb.GetNumber(&p.Currency))
	util.CheckError(pb.GetNumber(&p.SeatID))
	util.CheckError(pb.GetNumber(&p.CharID))
	util.CheckError(pb.GetNumber(&p.SkinID))
	util.CheckError(pb.GetNumber(&p.Radian))
	util.CheckError(pb.GetNumber(&p.BornTime))
	var buffCount uint8
	p.Buffs = make([]PtBuff, buffCount)
	for j := uint8(0); j < buffCount; j++ {
		util.CheckError(pb.GetNumber(&p.Buffs[j].BuffType))
		util.CheckError(pb.GetNumber(&p.Buffs[j].Data))
	}
}

type C2SHitFish struct {
	BulletSerial uint32
	OriginID uint32
	FishSerial uint32
	LocalTime float64
}

func (p *C2SHitFish) Bytes() []byte {
	var pb protocol.Protocol
	pb.SetCmd(HitFishCode)
	pb.AppendNumber(p)
	return pb.Bytes()
}

type dropItem struct {
	ModelID uint32
	Num uint32
}

type deadFish struct{
	Serial uint32
	IsDead uint8
	Items  []dropItem
}

type S2CHitFish struct {
	CharID     uint32
	SeatID     uint8
	Serial     uint32
	OriginID   uint32
	FishSerial uint32
	Multiple   uint32
	Currency   float64
	DeadFish   []deadFish
	ClientTime float64
}

func (p *S2CHitFish) Parse(pb *protocol.Protocol) {
	util.CheckError(pb.GetNumber(&p.CharID))
	util.CheckError(pb.GetNumber(&p.SeatID))
	util.CheckError(pb.GetNumber(&p.Serial))
	util.CheckError(pb.GetNumber(&p.OriginID))
	util.CheckError(pb.GetNumber(&p.FishSerial))
	util.CheckError(pb.GetNumber(&p.Multiple))
	util.CheckError(pb.GetNumber(&p.Currency))
	var count, dropCount uint8
	util.CheckError(pb.GetNumber(&count))
	p.DeadFish = make([]deadFish, count)
	for i := uint8(0); i < count; i++ {
		tempFish := &p.DeadFish[i]
		util.CheckError(pb.GetNumber(&tempFish.Serial))
		util.CheckError(pb.GetNumber(&tempFish.IsDead))
		util.CheckError(pb.GetNumber(&dropCount))
		tempFish.Items = make([]dropItem, dropCount)
		for j := uint8(0); j < dropCount; j++ {
			util.CheckError(pb.GetNumber(&tempFish.Items[j].ModelID))
			util.CheckError(pb.GetNumber(&tempFish.Items[j].Num))
		}
	}
	util.CheckError(pb.GetNumber(&p.ClientTime))
}

type C2STransmitActivity struct {
	Activity string
}

func (p *C2STransmitActivity) Bytes() []byte {
	var pb protocol.Protocol
	pb.SetCmd(TransmitCode)
	pb.AppendStringUint8(p.Activity)
	return pb.Bytes()
}

type S2CSyncBoom struct {
	Status     uint8
	LeftTime   uint32
}

func (p *S2CSyncBoom) Parse(pb *protocol.Protocol) {
	util.CheckError(pb.GetNumber(p))
}

type S2CGenerateFish struct {
	FishList []PtFish
}

func (p *S2CGenerateFish) Parse(pb *protocol.Protocol) {
	var count uint16
	util.CheckError(pb.GetNumber(&count))
	p.FishList = make([]PtFish, count)
	for i:=uint16(0); i < count; i++ {
		f := &p.FishList[i]
		util.CheckError(pb.GetNumber(&f.Serial))
		util.CheckError(pb.GetNumber(&f.KindID))
		util.CheckError(pb.GetNumber(&f.PathID))
		util.CheckError(pb.GetNumber(&f.Speed))
		util.CheckError(pb.GetNumber(&f.OffsetX))
		util.CheckError(pb.GetNumber(&f.OffsetY))
		util.CheckError(pb.GetNumber(&f.OffsetZ))
		util.CheckError(pb.GetNumber(&f.BornTime))
	}
}

type C2SDrawReadPacket struct {
	Grade uint8
}

func (p *C2SDrawReadPacket) Bytes() []byte {
	var pb protocol.Protocol
	pb.SetCmd(DrawReadPacketCode)
	pb.AppendNumber(p)
	return pb.Bytes()
}

type S2CRedPacketInfo struct {
	Chip uint64
	IsNewPlayer bool
}

func (p *S2CRedPacketInfo) Parse(pb *protocol.Protocol) {
	info, err := pb.GetStringUint16()
	util.CheckError(err)
	if len(info) == 0 {
		return
	}
	jv := make(map[string]interface{})
	err = json.Unmarshal([]byte(info), &jv)
	util.CheckError(err)
	p.Chip = uint64(jv["chipin"].(float64))
	p.IsNewPlayer = uint8(jv["np"].(float64)) == 1
}

type S2CPoseidonStatus struct {
	Status            uint8
	CurrLoopEndTime   uint32
	NextLoopStartTime uint32
	StartTime         string
	EndTime           string
	PlayAnimation     uint8
}

func (p *S2CPoseidonStatus) Parse(pb *protocol.Protocol) {
	var err error
	util.CheckError(pb.GetNumber(&p.Status))
	util.CheckError(pb.GetNumber(&p.CurrLoopEndTime))
	util.CheckError(pb.GetNumber(&p.NextLoopStartTime))
	p.StartTime, err = pb.GetStringUint8()
	util.CheckError(err)
	p.EndTime, err = pb.GetStringUint8()
	util.CheckError(err)
	util.CheckError(pb.GetNumber(&p.PlayAnimation))
}

type C2SHitPoseidon struct {
	BulletSerial uint32
	OriginID     uint32
	LocalTime    float64
}

func (p *C2SHitPoseidon) Bytes() []byte {
	var pb protocol.Protocol
	pb.SetCmd(HitPoseidonCode)
	pb.AppendNumber(p)
	return pb.Bytes()
}

type S2CHitPoseidon struct {
	CharID uint32
	SeatID uint8
	Serial uint32
	OriginID uint32
	Cost uint32
	Currency float64
	LocalTime float64
}

func (p *S2CHitPoseidon) Parse(pb *protocol.Protocol) {
	util.CheckError(pb.GetNumber(p))
}

type C2SSwitchCaliber struct {
	Caliber uint32
}

func (p *C2SSwitchCaliber) Bytes() []byte {
	var pb protocol.Protocol
	pb.SetCmd(SwitchCaliberCode)
	pb.AppendNumber(p)
	return pb.Bytes()
}

type S2CSwitchCaliber struct {
	CharID uint32
	SeatID uint8
	Caliber uint32
}

func (p *S2CSwitchCaliber) Parse(pb *protocol.Protocol) {
	util.CheckError(pb.GetNumber(p))
}

type C2SLaunchMissile struct {
	MissileID  uint32
	TargetFish uint32
}

func (p *C2SLaunchMissile) Bytes() []byte {
	var pb protocol.Protocol
	pb.SetCmd(LaunchMissileCode)
	pb.AppendNumber(p)
	return pb.Bytes()
}

type S2CLaunchMissile struct {
	Result uint8
	CharID uint32
	ModelID uint32
	TargetFish uint32
	RewardModelID uint32
	RewardNum uint32
	Currency float64
}

func (p *S2CLaunchMissile) Parse(pb *protocol.Protocol) {
	util.CheckError(pb.GetNumber(p))
}