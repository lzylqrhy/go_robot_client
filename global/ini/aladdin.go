package ini

import (
	"github/go-robot/util"
	"gopkg.in/ini.v1"
	"log"
)

var AladdinSetting struct{
	ServerAddr string
	RoomID     uint
	Line uint
	Chip uint
}

func getAladdinConfig(){
	path := "./configs/aladdin.ini"
	conf, err := ini.Load(path)
	util.CheckError(err)

	section := conf.Section("aladdin")
	if section == nil {
		log.Fatalf("fruit setting is not existed in %s\n", path)
	}
	GameCommonSetting.Frame = getOptionUInt(section, "frame", 1)
	AladdinSetting.ServerAddr = section.Key("server_addr").String()
	if section.HasKey("room_id") {
		AladdinSetting.RoomID, err = section.Key("room_id").Uint()
		util.CheckError(err)
	} else {
		log.Fatalf("fruit setting need room_id option in %s\n", path)
	}
	AladdinSetting.Line = getOptionUInt(section, "line", 50)
	AladdinSetting.Chip = getOptionUInt(section, "chip", 1)
}

