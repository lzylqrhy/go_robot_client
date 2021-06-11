package global

import (
	"github/go-robot/util"
	"gopkg.in/ini.v1"
	"log"
	"strconv"
	"strings"
)

var FishSetting struct {
	ServerAddr string
	RoomID     uint
	CaptureFishType map[uint32]struct{}
	HitPoseidon uint
	Caliber uint
}

var FishTestDataSetting struct{
	Items     map[uint]uint
	CaliberLV uint
}

func getFishConfig(){
	path := "./configs/fish.ini"
	conf, err := ini.Load(path)
	util.CheckError(err)

	section := conf.Section("fish")
	if section == nil {
		log.Fatalf("fish setting is not existed in %s\n", path)
		return
	}
	GameCommonSetting.Frame = getOptionUInt(section, "frame", 5)
	FishSetting.ServerAddr = section.Key("server_addr").String()
	FishSetting.RoomID = getOptionUInt(section, "room_id", 0)

	FishSetting.CaptureFishType = make(map[uint32]struct{})
	if section.HasKey("capturing_fish_type") {
		tempValue, err := section.Key("capturing_fish_type").StrictInt64s(",")
		util.CheckError(err)
		var empty struct{}
		for _, v := range tempValue {
			FishSetting.CaptureFishType[uint32(v)] = empty
		}
	}

	FishSetting.HitPoseidon = getOptionUInt(section, "hit_poseidon", 0)
	FishSetting.Caliber = getOptionUInt(section, "caliber", 0)

	// db
	getFishDBSetting(conf, "db_user", &GameCommonSetting.UserDB)
	getFishDBSetting(conf, "db_data", &GameCommonSetting.DataDB)
	// data
	getFishInitData(conf)
}

func getFishDBSetting(conf *ini.File, secKey string, setting *DBSetting)  {
	section := conf.Section(secKey)
	if section == nil {
		log.Printf("%s setting is not existed in fish.ini", secKey)
		return
	}
	setting.Account = section.Key("account").String()
	setting.Password = section.Key("password").String()
	setting.Database = section.Key("database").String()
	setting.Address = section.Key("addr").String()
	setting.Port = getOptionUInt(section,"port", 0)
	setting.IsUsable = true
}

func getFishInitData(conf *ini.File)  {
	section := conf.Section("test_data")
	if section == nil {
		log.Printf("%s setting is not existed in fish.ini", "test_data")
		return
	}
	FishTestDataSetting.Items = make(map[uint]uint)
	items := section.Key("items").String()
	itemsSlice := strings.Split(items, ",")
	for _, item := range itemsSlice {
		ivs := strings.Split(item, ":")
		if len(ivs) != 2 {
			continue
		}
		id, err := strconv.Atoi(ivs[0])
		if err != nil {
			log.Fatalln(err)
		}
		value, err := strconv.Atoi(ivs[1])
		if err != nil {
			log.Fatalln(err)
		}
		FishTestDataSetting.Items[uint(id)] = uint(value)
	}
	FishTestDataSetting.CaliberLV = getOptionUInt(section,"caliber_lv", 32)
}