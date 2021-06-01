package games

import (
	"github/go-robot/common"
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
	}
	log.Panic("game id is not undefined")
	return nil
}