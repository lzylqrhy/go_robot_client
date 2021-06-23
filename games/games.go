/**
 获取各游戏的机器人对象，并设置初始数据
 created by lzy
*/
package games

import (
	"github/go-robot/common"
	"github/go-robot/core"
	"github/go-robot/core/mynet"
	"github/go-robot/games/aladdin"
	"github/go-robot/games/fish"
	"github/go-robot/games/fruit"
	"github/go-robot/global"
	"github/go-robot/global/ini"
	"log"
)

func NewClient(index uint, pd *common.PlatformData, dialer mynet.MyDialer) core.RobotClient {
	switch ini.MainSetting.GameID {
	case global.FishGame:
		return fish.NewClient(index, pd, dialer)
	case global.FruitGame:
		return fruit.NewClient(index, pd, dialer)
	case global.AladdinGame:
		return aladdin.NewClient(index, pd, dialer)
	}
	log.Panic("game id is not undefined")
	return nil
}

func SetRobotTestData(pds []*common.PlatformData)  {
	switch ini.MainSetting.GameID {
	case global.FishGame:
		fish.RunTestData(pds)
	default:
		log.Printf("game %d have no test data", ini.MainSetting.GameID)
	}
	log.Printf("game %d set robot data completed", ini.MainSetting.GameID)
}
