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

var GameCommonSetting struct{
	Frame uint
}

var FishSetting struct {
	ServerAddr string
	RoomID     uint
	CaptureFishType map[uint32]struct{}
	HitPoseidon uint
}

var FruitSetting struct{
	ServerAddr string
	RoomID     uint
	Line uint
	Chip uint
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
	conf, err := ini.Load("./configs/game.ini")
	util.CheckError(err)
	GameCommonSetting.Frame, err = conf.Section("common").Key("frame").Uint()
	util.CheckError(err)
	if 0 == GameCommonSetting.Frame {
		log.Fatalln("game.common.frame is > 0")
	}
	switch MainSetting.GameID {
	case FishGame:
		getFishConfig(conf)
	case FruitGame:
		getFruitConfig(conf)
	case AladdinGame:
	default:
		log.Println("load ./configs/game.ini completed")
		break
	}
	log.Println("load ./configs/game.ini completed")
}

func getFishConfig(conf *ini.File){
	section := conf.Section("fish")
	if section == nil {
		log.Fatalln("fish setting is not existed in game.ini")
	}
	var err error
	FishSetting.ServerAddr = section.Key("server_addr").String()
	if section.HasKey("room_id") {
		FishSetting.RoomID, err = section.Key("room_id").Uint()
		util.CheckError(err)
	}
	FishSetting.CaptureFishType = make(map[uint32]struct{})
	if section.HasKey("capturing_fish_type") {
		tempValue, err := section.Key("capturing_fish_type").StrictInt64s(",")
		util.CheckError(err)
		var empty struct{}
		for _, v := range tempValue {
			FishSetting.CaptureFishType[uint32(v)] = empty
		}
	}
	if section.HasKey("hit_poseidon") {
		FishSetting.HitPoseidon, err = section.Key("hit_poseidon").Uint()
		util.CheckError(err)
	}
}

func getFruitConfig(conf *ini.File){
	section := conf.Section("fruit")
	if section == nil {
		log.Fatalln("fruit setting is not existed in game.ini")
	}
	var err error
	FruitSetting.ServerAddr = section.Key("server_addr").String()
	if section.HasKey("room_id") {
		FruitSetting.RoomID, err = section.Key("room_id").Uint()
		util.CheckError(err)
	} else {
		log.Fatalln("fruit setting need room_id option in game.ini")
	}
	FruitSetting.Line = getOptionUInt(section, "line", 9)
	FruitSetting.Chip = getOptionUInt(section, "chip", 1)
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