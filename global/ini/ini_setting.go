package ini

import (
	"github/go-robot/global"
	"github/go-robot/util"
	"gopkg.in/ini.v1"
	"log"
	"strconv"
)

// 主配置
var MainSetting struct{
	NetProtocol string
	RobotAddr string
	RobotStart  uint
	RobotNum    uint
	GameID      uint
	GameZone	string
}
// 数据库配置
type DBSetting struct {
	Account, Password, Database, Address string
	Port                                 uint
	IsUsable                             bool
}
// 通用配置
var GameCommonSetting struct{
	Frame          uint
	UserDB, DataDB DBSetting
}
// 加载ini配置
func LoadSetting()  {
	loadMainSetting()
	loadGameSetting()
}

func loadMainSetting() {
	conf, err := ini.Load("./configs/main.ini")
	util.CheckError(err)
	// 机器人
	robot := conf.Section("robot")
	if robot == nil {
		log.Panicln("must set main.robot section")
		return
	}
	MainSetting.NetProtocol = robot.Key("protocol").String()
	if "" == MainSetting.NetProtocol {
		log.Panicln("must set main.robot.protocol")
		return
	}
	MainSetting.RobotAddr = robot.Key("api_addr").String()
	if "" == MainSetting.RobotAddr {
		log.Panicln("must set main.robot.api_addr")
		return
	}
	MainSetting.RobotStart = getOptionUInt(robot, "start", 2)
	MainSetting.RobotNum = getOptionUInt(robot, "num", 1)
	MainSetting.GameID = getOptionUInt(robot, "game_id", 0)
	MainSetting.GameZone = strconv.Itoa(getOptionInt(robot, "game_zone", 0))
	log.Println("load ./configs/main.ini completed")
}

func loadGameSetting() {
	switch MainSetting.GameID {
	case global.FishGame:
		getFishConfig()
	case global.FruitGame:
		getFruitConfig()
	case global.AladdinGame:
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

func getDBSetting(conf *ini.File, secKey string, setting *DBSetting)  {
	section := conf.Section(secKey)
	if section == nil {
		log.Printf("%s setting is not existed in game[%d].ini", secKey, MainSetting.GameID)
		return
	}
	setting.Account = section.Key("account").String()
	setting.Password = section.Key("password").String()
	setting.Database = section.Key("database").String()
	setting.Address = section.Key("addr").String()
	setting.Port = getOptionUInt(section,"port", 0)
	if setting.Database != "" {
		setting.IsUsable = true
	}
}