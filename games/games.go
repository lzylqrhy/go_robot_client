package games

import (
	"context"
	"github/go-robot/common"
	"github/go-robot/games/aladdin"
	"github/go-robot/games/fish"
	"github/go-robot/games/fruit"
	"github/go-robot/global"
	myNet "github/go-robot/net"
	"log"
)

func NewClient(gameID uint, index uint, pd *common.PlatformData, dialer myNet.MyDialer) common.Client {
	switch gameID {
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

func SetRobotTestData(ctx context.Context, gameID uint, pID uint32)  {
	var isOK bool
	switch gameID {
	case global.FishGame:
		isOK = fish.SetRobotTestData(ctx, pID)
	default:
		log.Println("have no test data")
	}
	if !isOK {
		log.Fatalln("set robot data failed")
	} else {
		log.Println("set robot data successfully")
	}
}
