package games

import (
	"github/go-robot/common"
	"github/go-robot/games/fish"
	myNet "github/go-robot/net"
	"log"
)

const (
	FishGame = 1 + iota
)

func NewClient(gameID uint, index uint, pd *common.PlatformData, dialer myNet.MyDialer) common.Client {
	switch gameID {
	case FishGame:
		return fish.NewClient(index, pd, dialer)
	}
	log.Panic("game id is not undefined")
	return nil
}