package global

import (
	"github/go-robot/util"
	"gopkg.in/ini.v1"
	"log"
)

var MainSetting struct{
	NetProtocol string
	RobotAddr string
	RobotStart  uint
	RobotNum    uint
	GameID      uint
	GameZone	string
}

type DBSetting struct {
	Account, Password, Database, Address string
	Port                                 uint
	IsUsable                             bool
}

var GameCommonSetting struct{
	Frame          uint
	UserDB, DataDB DBSetting
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
		log.Panicln("must set main.server.protocol")
		return
	}
	// 机器人
	robot := conf.Section("robot")
	if robot != nil {
		MainSetting.RobotAddr = robot.Key("api_addr").String()
		if "" == MainSetting.RobotAddr {
			log.Panicln("must set main.robot.api_addr")
			return
		}
		MainSetting.RobotStart, err = robot.Key("start").Uint()
		util.CheckError(err)
		MainSetting.RobotNum, err = robot.Key("num").Uint()
		util.CheckError(err)
		MainSetting.GameID, err = robot.Key("game_id").Uint()
		util.CheckError(err)
		MainSetting.GameZone = robot.Key("game_zone").String()
		if "" == MainSetting.GameZone {
			log.Panicln("must set main.robot.game_zone")
			return
		}
	} else {
		log.Panicln("must set main.robot section")
		return
	}
	log.Println("load ./configs/main.ini completed")
}

func loadGameSetting() {
	switch MainSetting.GameID {
	case FishGame:
		getFishConfig()
	case FruitGame:
		getFruitConfig()
	case AladdinGame:
		getAladdinConfig()
	default:
		log.Printf("don't find ini of game id %d\n", MainSetting.GameID)
		break
	}
	log.Printf("load ini game id %d completed\n", MainSetting.GameID)
}

func getOptionInt(section *ini.Section, key string, def int) int {
	if section.HasKey(key) {
		v, err := section.Key(key).Int()
		util.CheckError(err)
		return v
	}
	return def
}

func getOptionUInt(section *ini.Section, key string, def uint) uint {
	if section.HasKey(key) {
		v, err := section.Key(key).Uint()
		util.CheckError(err)
		return v
	}
	return def
}
