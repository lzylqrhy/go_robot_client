package fish

type Buff struct {

}

type Fish struct {

}

type Bullet struct {

}

type Player struct {
	CharID uint32
	GameCurrency uint64
	SeatID uint8
	CannonID uint32
	Caliber uint32
	CaliberLV uint8
	Status uint16
	
}

type Pond struct {
	mapFish map[uint32]Fish
	mapBullet map[uint32]Bullet
	mapPlayer map[uint32]Player
}