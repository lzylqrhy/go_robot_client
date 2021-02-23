package global

import (
	"github/go-robot/util"
	"gopkg.in/ini.v1"
	"log"
)

var MainSetting struct{
	NetProtocol string
	RobotStart  uint
	RobotNum    uint
	GameID      uint
	GameZone	string
}

var FishSetting struct {
	ServerAddr string
	RoomID     uint
}

func LoadSetting()  {
	loadMainSetting()
	loadGameSetting()
}

func loadMainSetting() {
	conf, err := ini.Load("./configs/main.ini")
	util.CheckError(err)
	MainSetting.NetProtocol = conf.Section("server").Key("protocol").String()
	if "" == MainSetting.NetProtocol {
		panic("must set main.server.protocol")
		return
	}
	// 机器人数量
	robot := conf.Section("robot")
	if robot != nil {
		MainSetting.RobotStart, err = robot.Key("start").Uint()
		util.CheckError(err)
		MainSetting.RobotNum, err = robot.Key("num").Uint()
		util.CheckError(err)
		MainSetting.GameID, err = robot.Key("game_id").Uint()
		util.CheckError(err)
		MainSetting.GameZone = robot.Key("game_zone").String()
		if "" == MainSetting.GameZone {
			panic("must set main.robot.game_zone")
			return
		}
	} else {
		panic("must set main.robot section")
		return
	}
	log.Println("load ./configs/main.ini completed")
}

func loadGameSetting() {
	conf, err := ini.Load("./configs/game.ini")
	util.CheckError(err)
	section := conf.Section("fish")
	if section != nil {
		FishSetting.ServerAddr = section.Key("server_addr").String()
		if section.HasKey("room_id") {
			FishSetting.RoomID, err = section.Key("room_id").Uint()
			util.CheckError(err)
		}
	}
	log.Println("load ./configs/game.ini completed")
}