package games

import (
	"github/go-robot/common"
	"github/go-robot/games/aladdin"
	"github/go-robot/games/fish"
	"github/go-robot/games/fruit"
	"github/go-robot/global"
	myNet "github/go-robot/net"
	"log"
)

func NewClient(index uint, pd *common.PlatformData, dialer myNet.MyDialer) common.Client {
	switch global.MainSetting.GameID {
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
	switch global.MainSetting.GameID {
	case global.FishGame:
		fish.RunTestData(pds)
	default:
		log.Printf("game %d have no test data", global.MainSetting.GameID)
	}
	log.Printf("game %d set robot data completed", global.MainSetting.GameID)
}
