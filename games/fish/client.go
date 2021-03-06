package fish

import (
	"github/go-robot/common"
	"github/go-robot/core"
	"github/go-robot/core/mynet"
	"github/go-robot/core/protocol"
	"github/go-robot/global/ini"
	"github/go-robot/protocols"
	"log"
	"math"
	"math/rand"
	"time"
)

type FClient struct {
	common.ClientBase        // 公共数据
	pond              pond // 鱼池
	charID            uint32
	gameCurrency      uint64	// 游戏币
	seatID            uint8
	cannonID          uint32
	caliber           uint32
	caliberLV         uint8
	status            uint16           // 游戏状态（来自服务端）
	originSerial      uint32           // 子弹最新本地序号
	bulletCache       []bullet         // 子弹缓存
	rooms             []protocols.Room // 房间列表
	fireTime,hitTime  map[uint32]int64
	getInfoTime       int64
	poseidonStatus    uint8  //波塞冬游戏状态
}

func NewClient(index uint, pd *common.PlatformData, dialer mynet.MyDialer) core.RobotClient {
	c := new(FClient)
	c.Index = uint32(index)
	c.PtData = pd
	c.Dialer = dialer
	c.pond.Init()
	c.bulletCache = make([]bullet, 0, 20)
	c.fireTime = make(map[uint32]int64)
	c.hitTime = make(map[uint32]int64)
	log.Printf("new client: %v \n", c)
	log.Printf("new client: %v \n", *c.PtData)
	return c
}

// 获取服务端的毫秒
func (c *FClient) getServerTime() uint64 {
	if 0 == c.LocalTime {
		return 0
	}
	return c.SevTime + (uint64(time.Now().UnixNano() / 1e6) - c.LocalTime)
}

func (c *FClient)Update() {
	//log.Printf("client serial=%d update\n", c.serial)
	if !c.IsWorking {
		return
	}
	// 更新鱼坐标
	c.pond.mapFish.Update(c.getServerTime())
	// 处理缓存中的子弹
	c.cleanBulletCache()
	// 开火
	c.fire()
	// 发射导弹
	c.launchMissile()
}

func (c *FClient) cleanBulletCache() {
	if len(c.pond.mapFish) == 0 {
		return
	}
	isHit := false
	for _, b := range c.bulletCache {
		serial := c.getOneFish()
		if serial > 0 {
			isHit = true
			// 发送命中鱼
			var c2sHit = protocols.C2SHitFish{
				BulletSerial: b.Serial,
				OriginID: b.OriginID,
				FishSerial: serial,
			}
			c.SendPacket(c2sHit.Bytes())
			c.hitTime[b.Serial] = time.Now().UnixNano()
		}
	}
	if isHit {
		c.bulletCache = c.bulletCache[:0]
	}
}

func (c *FClient)OnConnected()  {
	var ping protocols.C2SSyncTime
	c.SendPacket(ping.Bytes())
	log.Printf("client index=%d, pid=%d connected\n", c.Index, c.PtData.PID)
}

func (c *FClient)OnDisconnected()  {
	log.Printf("client index=%d, pid=%d disconnected\n", c.Index, c.PtData.PID)
	c.IsWorking = false
}

func (c *FClient)ProcessProtocols(p *protocol.Protocol) bool {
	//log.Printf("cmd:0x%04x\n", p.Head.Cmd)
	isCommon, isOK := c.ProcessCommonProtocols(p)
	if !isCommon {
		switch p.Head.Cmd {
		case protocols.S2CLoginCode:
			isOK = c.processLogin(p)
		case protocols.EnterHallOrRoomCode:
			isOK = c.processEnterHallOrRoom(p)
		case protocols.ReadPacketInfoCode:
			isOK = c.processDrawRedPacket(p)
		case protocols.RoomListCode:
			isOK = c.processRoomList(p)
		case protocols.PlayerCode:
			isOK = c.processPlayerInfo(p)
		case protocols.FishEnterRoomCode:
			isOK = c.processEnterRoom(p)
		case protocols.SceneInfoCode:
			isOK = c.processSceneInfo(p)
		case protocols.PlayerSeatCode:
			isOK = c.processSeatsInfo(p)
		case protocols.FishListCode:
			isOK = c.processFishList(p)
		case protocols.BulletListCode:
			isOK = c.processBulletList(p)
		case protocols.FireCode:
			isOK = c.processFire(p)
		case protocols.HitFishCode:
			isOK = c.processHitFish(p)
		case protocols.SyncFishBoom:
			isOK = c.processSyncFishBoom(p)
		case protocols.GenerateFish:
			isOK = c.processGenerateFish(p)
		case protocols.PoseidonStatusCode:
			isOK = c.processPoseidonStatus(p)
		case protocols.HitPoseidonCode:
			isOK = c.processHitPoseidon(p)
		case protocols.SwitchCaliberCode:
			isOK = c.processSwitchCaliber(p)
		case protocols.LaunchMissileCode:
			isOK = c.processLaunchMissile(p)
		default:
			log.Printf("cmd:0x%04x don't be processed\n", p.Head.Cmd)
			isOK = true
		}
	}
	if !isOK {
		log.Printf("process cmd:0x%04x unsuccessfully, will exit\n", p.Head.Cmd)
	}
	return isOK
}

func (c *FClient) processLogin(p *protocol.Protocol) bool {
	var s2cLogin protocols.S2CLogin
	s2cLogin.Parse(p)
	if s2cLogin.Status == 1 {
		// 登录成功
		log.Printf("client index=%d, pid=%d login successfully\n", c.Index, c.PtData.PID)
		// 发送资源加载完成
		var c2sLoaded protocols.C2SResourceLoaded
		c.SendPacket(c2sLoaded.Bytes())
		return true
	}
	log.Printf("client index=%d, pid=%d login failed, status: %d\n", c.Index, c.PtData.PID, s2cLogin.Status)
	return false
}

func (c *FClient) processPlayerInfo(p *protocol.Protocol) bool {
	var s2cPlayer protocols.S2CPlayerInfo
	s2cPlayer.Parse(p)
	log.Printf("client index=%d, pid=%d get player info successfully, player=%v\n", c.Index, c.PtData.PID, s2cPlayer)
	c.charID = s2cPlayer.CharID
	return true
}

func (c *FClient) processEnterHallOrRoom(p *protocol.Protocol) bool {
	var s2cGo protocols.S2CEnterHallOrRoom
	s2cGo.Parse(p)
	// 如果有房间ID，则说明是断线重连，进入对应房间
	var roomID uint32
	if s2cGo.RoomID > 0 {
		roomID = s2cGo.RoomID
	}else {
		// 是否有指定房间
		if ini.FishSetting.RoomID > 0 {
			roomID = uint32(ini.FishSetting.RoomID)
		}else {
			// 找到可进入的房间
			validRoom := make([]uint32, 0, 3)
			for _, r := range c.rooms {
				// 目前此处缺少个人炮倍及游戏币数据
				//if c.gameCurrency < uint64(r.MinScore) ||
				//	r.MaxScore > 0 && c.gameCurrency > uint64(r.MaxScore) || c.caliber < r.MinScore{
				//	continue
				//}
				validRoom = append(validRoom, r.RoomID)
			}
			index := rand.Int31n(int32(len(validRoom)))
			roomID = validRoom[index]
		}
	}
	// 进入房间
	var c2sEnterRoom protocols.C2SFishEnterRoom
	c2sEnterRoom.RoomID = roomID
	c.SendPacket(c2sEnterRoom.Bytes())
	log.Printf("client index=%d, pid=%d player enter room %d\n", c.Index, c.PtData.PID, roomID)
	return true
}

func (c *FClient) processRoomList(p *protocol.Protocol) bool {
	var s2cRooms protocols.S2CRoomList
	s2cRooms.Parse(p)
	c.rooms = s2cRooms.Rooms
	log.Printf("client index=%d, pid=%d get room list:%v\n", c.Index, c.PtData.PID, c.rooms)
	return true
}

func (c *FClient) processEnterRoom(p *protocol.Protocol) bool {
	var s2cEnterRoom protocols.S2CFishEnterRoom
	s2cEnterRoom.Parse(p)
	if s2cEnterRoom.Result != 0 {
		// 进入失败
		log.Printf("client index=%d, pid=%d enter room=%d failed, result: %d\n",
			c.Index, c.PtData.PID, s2cEnterRoom.RoomID, s2cEnterRoom.Result)
		return false
	}
	c.RoomID = s2cEnterRoom.RoomID
	log.Printf("client index=%d, pid=%d player enter room=%d successfully\n", c.Index, c.PtData.PID, s2cEnterRoom.RoomID)
	// 请求场景信息
	var c2sSceneInfo protocols.C2SGetSceneInfo
	c.SendPacket(c2sSceneInfo.Bytes())
	c.getInfoTime = time.Now().UnixNano() / 1e6
	return true
}

func (c *FClient) processSceneInfo(p *protocol.Protocol) bool {
	var s2cSceneInfo protocols.S2CGetSceneInfo
	s2cSceneInfo.Parse(p)
	log.Printf("client index=%d, pid=%d get scene info successfully\n", c.Index, c.PtData.PID)
	c.SevTime = uint64(s2cSceneInfo.ServerTime)
	c.LocalTime = uint64(time.Now().UnixNano() / 1e6)
	return true
}

func (c *FClient) processSeatsInfo(p *protocol.Protocol) bool {
	var s2cSeats protocols.S2CSeatsInfo
	s2cSeats.Parse(p)
	log.Printf("client index=%d, pid=%d get seat list successfully\n", c.Index, c.PtData.PID)
	for _, p := range s2cSeats.Players {
		if p.CharID == c.charID {
			c.seatID = p.SeatID
			c.caliber = p.Caliber
			c.caliberLV = p.CaliberLV
			c.cannonID = p.CannonID
			c.gameCurrency = uint64(p.Currency)
			c.status = p.Status
			continue
		}
		c.pond.mapPlayer[p.CharID] = player{
			CharID: p.CharID,
			SeatID: p.SeatID,
			GameCurrency: uint64(p.Currency),
			CannonID: p.CannonID,
			Caliber: p.Caliber,
			CaliberLV: p.CaliberLV,
			Status: p.Status,
		}
	}
	// 切换炮倍
	c.switchCaliber()
	return true
}

func (c *FClient) processFishList(p *protocol.Protocol) bool {
	var s2cFish protocols.S2CFishList
	s2cFish.Parse(p)
	log.Printf("client index=%d, pid=%d get fish list successfully\n", c.Index, c.PtData.PID)
	for _, f := range s2cFish.FishList {
		c.pond.mapFish[f.Serial] = fish{
			Serial:   f.Serial,
			KindID:   f.KindID,
			PathID:   f.PathID,
			Speed:    f.Speed,
			OffsetX:  f.OffsetX,
			OffsetY:  f.OffsetY,
			OffsetZ:  f.OffsetZ,
			BornTime: f.BornTime,
			SwamTime: f.SwamTime,
		}
	}
	return true
}

func (c *FClient) processBulletList(p *protocol.Protocol) bool {
	var s2cBullet protocols.S2CBulletList
	s2cBullet.Parse(p)
	log.Printf("client index=%d, pid=%d get bullet list successfully\n", c.Index, c.PtData.PID)
	for _, b := range s2cBullet.BulletList {
		if b.CharID == c.charID {
			// 添加到子弹列表中
			c.bulletCache = append(c.bulletCache, bullet{
				Serial: b.Serial,
				OriginID: b.OriginID,
				SeatID: b.SeatID,
				CharID: b.CharID,
				SkinID: b.SkinID,
				Radian: b.Radian,
				BornTime: b.BornTime,
			})
			continue
		}
		//c.pond.mapBullet[b.Serial] = bullet{
		//	Serial: b.Serial,
		//	OriginID: b.OriginID,
		//	SeatID: b.SeatID,
		//	CharID: b.CharID,
		//	SkinID: b.SkinID,
		//	Radian: b.Radian,
		//	BornTime: b.BornTime,
		//}
	}
	c.IsWorking = true
	//log.Println("----------------------------c.getInfo cost", time.Now().UnixNano() / 1e6 - c.getInfoTime)
	return true
}

const CostWarning = 0

func (c *FClient) fire() {
	if !ini.FishSetting.DoFire {
		return
	}
	// 判断钱是否足够
	if c.gameCurrency < uint64(c.caliber) {
		c.Disconnect("client index=%d charid=%d has no enough coin, need=%d, cur=%d, will exit \n", c.Index, c.charID, c.caliber, c.gameCurrency)
		return
	}
	c.gameCurrency -= uint64(c.caliber)
	// 判断缓存的子弹是否达到上限
	if len(c.bulletCache) >= 20 {
		log.Printf("client index=%d charid=%d has max buttel count >= 20 \n", c.Index, c.charID)
		return
	}
	// 发射子弹
	c.originSerial++
	c2sFire := protocols.C2SFire{}
	c2sFire.OriginID = c.originSerial
	c2sFire.Radian = float32(rand.Int31n(120) + 30) / 180 * math.Pi
	if c.seatID < 2 {
		c2sFire.Radian *= -1
	}
	c.SendPacket(c2sFire.Bytes())
	c.fireTime[c2sFire.OriginID] = time.Now().UnixNano()
}

func (c *FClient) processFire(p *protocol.Protocol) bool {
	var s2cFire protocols.S2CFire
	s2cFire.Parse(p)
	// 暂时只管自己的子弹
	if s2cFire.CharID != c.charID {
		return true
	}
	costTime := (time.Now().UnixNano() - c.fireTime[s2cFire.OriginID]) / 1e6
	if costTime > CostWarning {
		log.Printf("client index=%d, pid=%d fire cost time %d\n", c.Index, c.PtData.PID, costTime)
	}
	delete(c.fireTime, s2cFire.OriginID)
	if 0 != s2cFire.Result {
		log.Printf("client index=%d, pid=%d fire failed, result=%d \n", c.Index, c.PtData.PID, s2cFire.Result)
		return true
	}

	// 更新游戏币
	c.gameCurrency = uint64(s2cFire.Currency)
	log.Printf("client index=%d, pid=%d fire successfully, money=%d\n", c.Index, c.PtData.PID, c.gameCurrency)

	// 如果波塞冬房间，且要攻击波塞冬，则在波塞冬出现期间只打波塞冬
	if c.canHitPoseidon() {
		// 发送命中波塞冬
		var c2sHit = protocols.C2SHitPoseidon{
			BulletSerial: s2cFire.Serial,
			OriginID: s2cFire.OriginID,
		}
		c.SendPacket(c2sHit.Bytes())
		return true
	}

	// 获取一条鱼
	serial := c.getOneFish()
	if serial > 0 {
		// 发送命中鱼
		var c2sHit = protocols.C2SHitFish{
			BulletSerial: s2cFire.Serial,
			OriginID: s2cFire.OriginID,
			FishSerial: serial,
		}
		c.SendPacket(c2sHit.Bytes())
		c.hitTime[s2cFire.Serial] = time.Now().UnixNano()
	}else {
		// 添加到子弹列表中
		c.bulletCache = append(c.bulletCache, bullet{
			Serial: s2cFire.Serial,
			OriginID: s2cFire.OriginID,
			SeatID: s2cFire.SeatID,
			CharID: s2cFire.CharID,
			SkinID: s2cFire.SkinID,
			Radian: s2cFire.Radian,
			BornTime: s2cFire.BornTime,
		})
	}
	return true
}

func (c *FClient) getOneFish() uint32 {
	if 0 == len(c.pond.mapFish) {
		log.Printf("client index=%d, pid=%d has no any fish\n", c.Index, c.PtData.PID)
		return 0
	}
	fishShow := make([]uint32, 0, 32) // 可绘制出来的鱼
	mapFishType := make(map[uint32]uint32) // 当前鱼池鱼类型
	checkType := len(ini.FishSetting.CaptureFishType) > 0
	for k, v := range c.pond.mapFish {
		if v.BornTime <= float64(c.getServerTime()) {
			if checkType {
				t := ConfMgr.getFishTypeByID(v.KindID)
				mapFishType[t]++
				if _, isExisting := ini.FishSetting.CaptureFishType[t]; !isExisting {
					continue
				}
				//log.Printf("client index=%d, pid=%d can capture fish=%d, type=%d\n", c.Index, c.PtData.PID, k, t)
			}
			fishShow = append(fishShow, k)
		}
	}
	count := len(fishShow)
	if 0 == count {
		log.Printf("client index=%d, pid=%d has no target fish, all fish count=%d, fish types is %v\n", c.Index, c.PtData.PID, len(c.pond.mapFish), mapFishType)
		return 0
	}
	log.Printf("client index=%d, pid=%d target fish count=%d\n", c.Index, c.PtData.PID, count)
	index := rand.Int31n(int32(count))
	return fishShow[index]
}

func (c *FClient) processHitFish(p *protocol.Protocol) bool {
	var s2cHit protocols.S2CHitFish
	s2cHit.Parse(p)
	for _, f := range s2cHit.DeadFish {
		if f.IsDead > 0 {
			delete(c.pond.mapFish, f.Serial)
			//log.Printf("client index=%d, pid=%d captured fish=%d\n", c.Index, c.PtData.PID, f.Serial)
		}
	}
	if s2cHit.CharID == c.charID {
		c.gameCurrency = uint64(s2cHit.Currency)
		costTime := (time.Now().UnixNano() - c.hitTime[s2cHit.Serial]) / 1e6
		if costTime > CostWarning {
			log.Printf("client index=%d, pid=%d hit cost time %d\n", c.Index, c.PtData.PID, costTime)
		}
		delete(c.hitTime, s2cHit.Serial)
	}
	return true
}

func (c *FClient) processGenerateFish(p *protocol.Protocol) bool {
	var s2cGen protocols.S2CGenerateFish
	s2cGen.Parse(p)
	for _, f := range s2cGen.FishList {
		c.pond.mapFish[f.Serial] = fish{
			Serial:   f.Serial,
			KindID:   f.KindID,
			PathID:   f.PathID,
			Speed:    f.Speed,
			OffsetX:  f.OffsetX,
			OffsetY:  f.OffsetY,
			OffsetZ:  f.OffsetZ,
			BornTime: f.BornTime,
			SwamTime: f.SwamTime,
		}
	}
	return true
}

func (c *FClient) processSyncFishBoom(p *protocol.Protocol) bool {
	var s2cBoom protocols.S2CSyncBoom
	s2cBoom.Parse(p)
	// 鱼潮开始和结束时，清空鱼
	if s2cBoom.Status == 1 || s2cBoom.Status == 3 {
		// 清空场景中的鱼
		c.pond.mapFish.Clear()
	}
	return true
}

func (c *FClient) processDrawRedPacket(p *protocol.Protocol) bool {
	var s2cRedPacket protocols.S2CRedPacketInfo
	s2cRedPacket.Parse(p)
	// 获取红包配置
	conf := ConfMgr.getTownDrawRedPacketByID(c.RoomID)
	if nil == conf {
		log.Printf("client index=%d, pid=%d draw red packet, room=%d has no config \n", c.Index, c.PtData.PID, c.RoomID)
		return true
	}
	grade := uint8(3)
	isOK := false
	if s2cRedPacket.IsNewPlayer {
		isOK = s2cRedPacket.Chip >= conf.mapNewPlayerGrade[grade]
	}else {
		isOK = s2cRedPacket.Chip >= conf.mapNormalGrade[grade]
	}
	if isOK {
		c2sDraw := protocols.C2SDrawReadPacket{}
		c2sDraw.Grade = grade
		c.SendPacket(c2sDraw.Bytes())
		log.Printf("client index=%d, pid=%d draw red packet, grade=%d, is new player = %v  \n", c.Index, c.PtData.PID, grade, s2cRedPacket.IsNewPlayer)
	}
	return true
}

func (c *FClient) isPoseidonRoom() bool {
	for _, r := range c.rooms {
		if r.RoomID == c.RoomID {
			if r.Type == 4 {
				return true
			}
		}
	}
	return false
}

func (c *FClient) canHitPoseidon() bool {
	if c.isPoseidonRoom() && c.poseidonStatus == 2 && ini.FishSetting.HitPoseidon == 1 {
		return true
	}
	return false
}

func (c *FClient) processPoseidonStatus(p *protocol.Protocol) bool {
	var s2cPoseidonStatus protocols.S2CPoseidonStatus
	s2cPoseidonStatus.Parse(p)
	c.poseidonStatus = s2cPoseidonStatus.Status
	log.Printf("client index=%d, pid=%d get poseidonStatus successfully! now status[%d]\n", c.Index, c.PtData.PID, c.poseidonStatus)
	return true
}

func (c *FClient) processHitPoseidon(p *protocol.Protocol) bool {
	var s2cHitPoseidon protocols.S2CHitPoseidon
	s2cHitPoseidon.Parse(p)
	if s2cHitPoseidon.CharID == c.charID {
		c.gameCurrency = uint64(s2cHitPoseidon.Currency)
		log.Printf("client index=%d, pid=%d hit poseidon, current money[%d]\n", c.Index, c.PtData.PID, c.gameCurrency)
	}
	return true
}

func (c *FClient)switchCaliber()  {
	if caliber := uint32(ini.FishSetting.Caliber); ini.FishSetting.Caliber > 0 && caliber != c.caliber {
		c2sSwitchCaliber := protocols.C2SSwitchCaliber{Caliber: caliber}
		c.SendPacket(c2sSwitchCaliber.Bytes())
	}
}

func (c *FClient) processSwitchCaliber(p *protocol.Protocol) bool {
	var s2cSwitchCaliber protocols.S2CSwitchCaliber
	s2cSwitchCaliber.Parse(p)
	if s2cSwitchCaliber.CharID == c.charID {
		c.caliber = s2cSwitchCaliber.Caliber
	}
	return true
}

func (c *FClient) launchMissile() {
	if ini.FishSetting.LaunchMode == 0 {
		return
	}
	modelID := uint(0)
	for _, id := range ini.FishSetting.Missiles {
		if c.Items[uint32(id)] > 0 {
			modelID = id
			c.Items[uint32(id)]--
			break
		}
	}
	if modelID == 0 {
		return
	}
	// 找一条黄金鱼
	fishID := c.getOneSpecifiedFish(BigFish)
	if fishID > 0 {
		c2sLaunch := protocols.C2SLaunchMissile{MissileID: uint32(modelID), TargetFish: fishID}
		c.SendPacket(c2sLaunch.Bytes())
		log.Printf("client index=%d, pid=%d has no any fish\n", c.Index, c.PtData.PID)
	}
}

func (c *FClient) getOneSpecifiedFish(dstType uint32) uint32 {
	if 0 == len(c.pond.mapFish) {
		log.Printf("client index=%d, pid=%d has no any fish\n", c.Index, c.PtData.PID)
		return 0
	}
	fishShow := make([]uint32, 0, 4) // 可绘制出来的鱼
	mapFishType := make(map[uint32]uint32) // 当前鱼池鱼类型
	for k, v := range c.pond.mapFish {
		if v.BornTime <= float64(c.getServerTime()) {
			t := ConfMgr.getFishTypeByID(v.KindID)
			mapFishType[t]++
			if t != dstType {
				continue
			}
			fishShow = append(fishShow, k)
		}
	}
	count := len(fishShow)
	if 0 == count {
		log.Printf("client index=%d, pid=%d has no target fish, specified type %d, all fish count=%d, types is %v\n",
			c.Index, c.PtData.PID, dstType, len(c.pond.mapFish), mapFishType)
		return 0
	}
	//log.Printf("client index=%d, pid=%d target fish count=%d\n", c.Index, c.PtData.PID, count)
	index := rand.Int31n(int32(count))
	return fishShow[index]
}

func (c *FClient) processLaunchMissile(p *protocol.Protocol) bool {
	var s2cMissile protocols.S2CLaunchMissile
	s2cMissile.Parse(p)
	// 判断角色, 不是自己，直接返回
	if s2cMissile.CharID != c.charID {
		return true
	}
	// 失败，返回
	if s2cMissile.Result > 0 {
		log.Printf("client index=%d, pid=%d launch missile failed, result=%d\n", c.Index, c.PtData.PID, s2cMissile.Result)
		return true
	}
	// 奖励
	//c.Items[s2cMissile.ModelID]--
	c.Items[s2cMissile.RewardModelID] += uint64(s2cMissile.RewardNum)
	c.gameCurrency = uint64(s2cMissile.Currency)
	log.Printf("client index=%d, pid=%d launch missile successfully, model id=%d, left=%d, current money=%d\n",
		c.Index, c.PtData.PID, s2cMissile.ModelID, c.Items[s2cMissile.ModelID], c.gameCurrency)
	if ini.FishSetting.LaunchMode == 2 {
		// 遍历道具列表，只要有一种道具数量不为0，就继续
		for _, id := range ini.FishSetting.Missiles {
			if c.Items[uint32(id)] > 0 {
				return true
			}
		}
		c.Disconnect("client index=%d, pid=%d, char id=%d has no enough missiles, will exit \n", c.Index, c.PtData.PID, c.charID)
	}
	return true
}
