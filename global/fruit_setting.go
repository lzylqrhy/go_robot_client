package global

import (
	"github/go-robot/util"
	"gopkg.in/ini.v1"
	"log"
)

var FruitSetting struct{
	ServerAddr string
	RoomID     uint
	Line uint
	Chip uint
}

func getFruitConfig(){
	path := "./configs/fruit.ini"
	conf, err := ini.Load(path)
	util.CheckError(err)

	section := conf.Section("fruit")
	if section == nil {
		log.Fatalf("fruit setting is not existed in %s\n", path)
	}
	GameCommonSetting.Frame = getOptionUInt(section, "frame", 1)
	FruitSetting.ServerAddr = section.Key("server_addr").String()
	if section.HasKey("room_id") {
		FruitSetting.RoomID, err = section.Key("room_id").Uint()
		util.CheckError(err)
	} else {
		log.Fatalf("fruit setting need room_id option in %s\n", path)
	}
	FruitSetting.Line = getOptionUInt(section, "line", 9)
	FruitSetting.Chip = getOptionUInt(section, "chip", 1)
}

