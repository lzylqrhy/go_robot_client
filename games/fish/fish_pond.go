package fish

type buff struct {
	BuffType uint32
	Data 	 uint32
	Duration uint32
}

type fish struct {
	Serial   uint32
	FishID   uint32
	PathID   uint32
	Speed    uint32
	OffsetX  uint32
	OffsetY  uint32
	OffsetZ  uint32
	BornTime uint64
	SwamTime uint32
}

type bullet struct {
	Serial   uint32
	OriginID uint32
	SeatID   uint8
	CharID   uint32
	SkinID   uint32
	Radian   float32
	BornTime uint64
	Buffs    []buff
}

type player struct {
	CharID uint32
	GameCurrency uint64
	SeatID uint8
	CannonID uint32
	Caliber uint32
	CaliberLV uint8
	Status uint16
}

type fishManager map[uint32]fish

type pond struct {
	mapFish fishManager
	mapBullet map[uint32]bullet
	mapPlayer map[uint32]player
	buff []buff
}

func (p *pond) Init() {
	p.mapFish = make(fishManager, 32)
	p.mapBullet = make(map[uint32]bullet)
	p.mapPlayer = make(map[uint32]player, 3)
}

func (mgr *fishManager)Clear(){
	*mgr = make(fishManager, 32)
}
