package fish

import (
	"github/go-robot/common"
	myNet "github/go-robot/net"
	"github/go-robot/protocols"
	"log"
	"math/rand"
	"time"
)

type FClient struct {
	common.ClientBase      // 公共数据
	pond              pond // 鱼池
	charID            uint32
	gameCurrency      uint64	// 游戏币
	seatID            uint8
	cannonID          uint32
	caliber           uint32
	caliberLV         uint8
	status            uint16   // 游戏状态（来自服务端）
	isWork            bool     // 是否工作
	originSerial      uint32   // 子弹最新本地序号
	bulletCache       []bullet // 子弹缓存
}

func NewClient(index uint, pd *common.PlatformData, dialer myNet.MyDialer) common.Client {
	c := new(FClient)
	c.Index = uint32(index)
	c.PtData = pd
	c.Dialer = dialer
	c.pond.Init()
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
	if !c.isWork {
		return
	}
	// 更新鱼坐标
	c.pond.mapFish.Update(c.getServerTime())
	// 处理缓存中的子弹
	c.cleanBulletCache()
	// 开火
	c.fire()
}

func (c *FClient) cleanBulletCache() {
	if len(c.pond.mapFish) == 0 {
		return
	}
	for _, b := range c.bulletCache {
		serial := c.getOneFish()
		if serial > 0 {
			// 发送命中鱼
			var c2sHit = protocols.C2SHitFish{
				BulletSerial: b.Serial,
				OriginID: b.OriginID,
				FishSerial: serial,
			}
			c.SendPacket(c2sHit.Bytes())
		}
	}
	c.bulletCache = c.bulletCache[:0]
}

func (c *FClient)OnConnected()  {
	var ping protocols.C2SSyncTime
	c.SendPacket(ping.Bytes())
	log.Printf("client index=%d connected\n", c.Index)
}

func (c *FClient)OnDisconnected()  {
	log.Printf("client index=%d disconnected\n", c.Index)
	c.isWork = false
}

func (c *FClient)ProcessProtocols(p *protocols.Protocol) bool {
	//log.Printf("cmd:0x%04x\n", p.Head.Cmd)
	isCommon, isOK := c.ProcessCommonProtocols(p)
	if isCommon {
		return isOK
	}
	switch p.Head.Cmd {
	case protocols.S2CLoginCode:
		return c.processLogin(p)
	case protocols.PlayerCode:
		return c.processPlayerInfo(p)
	case protocols.EnterRoomCode:
		return c.processEnterRoom(p)
	case protocols.SceneInfoCode:
		return c.processSceneInfo(p)
	case protocols.PlayerSeatCode:
		return c.processSeatsInfo(p)
	case protocols.FishListCode:
		return c.processFishList(p)
	case protocols.BulletListCode:
		return c.processBulletList(p)
	case protocols.FireCode:
		return c.processFire(p)
	case protocols.HitFishCode:
		return c.processHitFish(p)
	case protocols.SyncFishBoom:
		return c.processSyncFishBoom(p)
	case protocols.GenerateFish:
		return c.processGenerateFish(p)
	}
	return true
}

func (c *FClient) processLogin(p *protocols.Protocol) bool {
	var s2cLogin protocols.S2CLogin
	s2cLogin.Parse(p)
	if s2cLogin.Status == 1 {
		// 登录成功
		log.Printf("client index=%d, pid=%d login successfully\n", c.Index, c.PtData.PID)
		// 发送资源加载完成
		var c2sLoaded protocols.C2SResourceLoaded
		c.SendPacket(c2sLoaded.Bytes())
		// 进入房间
		var c2sEnterRoom protocols.C2SEnterRoom
		c2sEnterRoom.RoomID = 20001
		c.SendPacket(c2sEnterRoom.Bytes())
		return true
	}
	log.Printf("client index=%d, pid=%d login failed, status: %d\n", c.Index, c.PtData.PID, s2cLogin.Status)
	return false
}

func (c *FClient) processPlayerInfo(p *protocols.Protocol) bool {
	var s2cPlayer protocols.S2CPlayerInfo
	s2cPlayer.Parse(p)
	log.Printf("client index=%d, pid=%d get player info successfully, player=%v\n", c.Index, c.PtData.PID, s2cPlayer)
	c.charID = s2cPlayer.CharID
	return true
}

func (c *FClient) processEnterRoom(p *protocols.Protocol) bool {
	var s2cEnterRoom protocols.S2CEnterRoom
	s2cEnterRoom.Parse(p)
	if s2cEnterRoom.Result != 0 {
		// 进入失败
		log.Printf("client index=%d, pid=%d enter room=%d failed, result: %d\n",
			c.Index, c.PtData.PID, s2cEnterRoom.RoomID, s2cEnterRoom.Result)
		return false
	}
	log.Printf("client index=%d, pid=%d enter room=%d successfully\n", c.Index, c.PtData.PID, s2cEnterRoom.RoomID)
	// 请求场景信息
	var c2sSceneInfo protocols.C2SGetSceneInfo
	c.SendPacket(c2sSceneInfo.Bytes())
	//// 转发
	//var c2sTransmit protocols.C2STransmitActivity
	//c2sTransmit.Activity = "hello"
	//c.sendPacket(c2sTransmit.Bytes())
	return true
}

func (c *FClient) processSceneInfo(p *protocols.Protocol) bool {
	var s2cSceneInfo protocols.S2CGetSceneInfo
	s2cSceneInfo.Parse(p)
	log.Printf("client index=%d, pid=%d get scene info successfully\n", c.Index, c.PtData.PID)
	c.SevTime = uint64(s2cSceneInfo.ServerTime)
	c.LocalTime = uint64(time.Now().UnixNano() / 1e6)
	return true
}

func (c *FClient) processSeatsInfo(p *protocols.Protocol) bool {
	var s2cSeats protocols.S2CSeatsInfo
	s2cSeats.Parse(p)
	log.Printf("client index=%d, pid=%d get seat list successfully\n", c.Index, c.PtData.PID)
	for _, p := range s2cSeats.Players {
		if p.CharID == c.charID {
			c.seatID = p.SeatID
			c.caliber = p.Caliber
			c.caliberLV = p.CaliberLV
			c.cannonID = p.CannonID
			c.gameCurrency = p.Currency
			c.status = p.Status
			continue
		}
		c.pond.mapPlayer[p.CharID] = player{
			CharID: p.CharID,
			SeatID: p.SeatID,
			GameCurrency: p.Currency,
			CannonID: p.CannonID,
			Caliber: p.Caliber,
			CaliberLV: p.CaliberLV,
			Status: p.Status,
		}
	}
	return true
}

func (c *FClient) processFishList(p *protocols.Protocol) bool {
	var s2cFish protocols.S2CFishList
	s2cFish.Parse(p)
	log.Printf("client index=%d, pid=%d get fish list successfully\n", c.Index, c.PtData.PID)
	for _, f := range s2cFish.FishList {
		c.pond.mapFish[f.Serial] = fish{
			Serial: f.Serial,
			FishID: f.FishID,
			PathID: f.PathID,
			Speed: f.Speed,
			OffsetX: f.OffsetX,
			OffsetY: f.OffsetY,
			OffsetZ: f.OffsetZ,
			BornTime: f.BornTime,
			SwamTime: f.SwamTime,
		}
	}
	return true
}

func (c *FClient) processBulletList(p *protocols.Protocol) bool {
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
		c.pond.mapBullet[b.Serial] = bullet{
			Serial: b.Serial,
			OriginID: b.OriginID,
			SeatID: b.SeatID,
			CharID: b.CharID,
			SkinID: b.SkinID,
			Radian: b.Radian,
			BornTime: b.BornTime,
		}
	}
	c.isWork = true
	return true
}

func (c *FClient) fire() {
	// 判断钱是否足够
	if c.gameCurrency < uint64(c.caliber) {
		log.Printf("client %d has no enough coin, need=%d, cur=%d, will exit \n", c.charID, c.caliber, c.gameCurrency)
		c.Dialer.Disconnect()
		c.isWork = false
		return
	}
	// 判断缓存的子弹是否达到上限
	if len(c.bulletCache) >= 20 {
		return
	}
	// 发射子弹
	c.originSerial++
	c2sFire := protocols.C2SFire{}
	c2sFire.OriginID = c.originSerial
	c2sFire.Radian = 0.8
	c.SendPacket(c2sFire.Bytes())
}

func (c *FClient) processFire(p *protocols.Protocol) bool {
	var s2cFire protocols.S2CFire
	s2cFire.Parse(p)
	if 0 != s2cFire.Result {
		log.Printf("client index=%d, pid=%d fire failed, result=%d \n", c.Index, c.PtData.PID, s2cFire.Result)
		return true
	}
	if s2cFire.CharID != c.charID {
		// 别人的子弹不管
		return true
	}
	//log.Printf("client index=%d, pid=%d fire successfully\n", c.Index, c.PtData.PID)
	// 更新游戏币
	c.gameCurrency = s2cFire.Currency
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
	count := len(c.pond.mapFish)
	if count == 0 {
		return 0
	}
	index := rand.Int31n(int32(count)) + 1
	for k := range c.pond.mapFish {
		index--
		if 0 == index {
			return k
		}
	}
	return 0
}

func (c *FClient) processHitFish(p *protocols.Protocol) bool {
	var s2cHit protocols.S2CHitFish
	s2cHit.Parse(p)
	for _, f := range s2cHit.DeadFish {
		if f.IsDead > 0 {
			delete(c.pond.mapFish, f.Serial)
			log.Printf("client index=%d, pid=%d captured fish\n", c.Index, c.PtData.PID)
		}
	}
	if s2cHit.CharID == c.charID {
		c.gameCurrency = s2cHit.Currency
	}
	return true
}

func (c *FClient) processGenerateFish(p *protocols.Protocol) bool {
	var s2cGen protocols.S2CGenerateFish
	s2cGen.Parse(p)
	for _, f := range s2cGen.FishList {
		c.pond.mapFish[f.Serial] = fish{
			Serial: f.Serial,
			FishID: f.FishID,
			PathID: f.PathID,
			Speed: f.Speed,
			OffsetX: f.OffsetX,
			OffsetY: f.OffsetY,
			OffsetZ: f.OffsetZ,
			BornTime: f.BornTime,
			SwamTime: f.SwamTime,
		}
	}
	return true
}

func (c *FClient) processSyncFishBoom(p *protocols.Protocol) bool {
	var s2cBoom protocols.S2CSyncBoom
	s2cBoom.Parse(p)
	// 鱼潮开始和结束时，清空鱼
	if s2cBoom.Status == 1 || s2cBoom.Status == 3 {
		// 清空场景中的鱼
		c.pond.mapFish.Clear()
	}
	return true
}