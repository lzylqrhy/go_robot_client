package ini

import (
	"github/go-robot/global"
	"github/go-robot/util"
	"gopkg.in/ini.v1"
	"log"
	"strconv"
	"strings"
)

var FishSetting struct {
	ServerAddr string
	RoomID     uint
	DoFire	bool
	CaptureFishType map[uint32]struct{}
	HitPoseidon uint
	Caliber uint
	LaunchMode uint
	Missiles []uint
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
	// 开火
	getFireSetting(section)
	// 波塞冬
	FishSetting.HitPoseidon = getOptionUInt(section, "hit_poseidon", 0)
	// 发射导弹
	getMissileSetting(section)
	// db
	getDBSetting(conf, "db_user", &GameCommonSetting.UserDB)
	getDBSetting(conf, "db_data", &GameCommonSetting.DataDB)
	// data
	getFishInitData(conf)
}

func getFireSetting(section *ini.Section)  {
	FishSetting.DoFire = getOptionUInt(section, "do_fire", 0) == 1
	FishSetting.CaptureFishType = make(map[uint32]struct{})
	if section.HasKey("capturing_fish_type") {
		tempValue, err := section.Key("capturing_fish_type").StrictInt64s(",")
		util.CheckError(err)
		var empty struct{}
		for _, v := range tempValue {
			FishSetting.CaptureFishType[uint32(v)] = empty
		}
	}
	FishSetting.Caliber = getOptionUInt(section, "caliber", 0)
}

func getMissileSetting(section *ini.Section) {
	FishSetting.LaunchMode = getOptionUInt(section, "launch_mode", 0)
	FishSetting.Missiles = make([]uint, 0)
	missileStr := section.Key("missile").String()
	if FishSetting.LaunchMode > 0 {
		temp := strings.Split(missileStr, ",")
		for _, ids := range temp {
			if id, err := strconv.Atoi(ids); err == nil {
				FishSetting.Missiles = append(FishSetting.Missiles, uint(id))
			}
		}
		if len(FishSetting.Missiles) == 0 {
			FishSetting.Missiles = append(FishSetting.Missiles, global.ItemFishBlackMissile, global.ItemFishBronzeMissile,
				global.ItemFishSilverMissile, global.ItemFishGoldMissile, global.ItemFishPlatinumMissile, global.ItemFishKingMissile)
			log.Println("current setting is launching all of missiles")
		}
	}
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