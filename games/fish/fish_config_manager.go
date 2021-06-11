package fish

import (
	"encoding/json"
	"errors"
	"github/go-robot/util"
	"io/ioutil"
	"log"
	"os"
)

// 路径配置
type pathConfig struct {
	ID uint32
	SwimDuration uint32
	StayDuration uint32
}

// 房间开红包配置
type roomDrawRedPacket struct {
	townID uint32
	mapNormalGrade, mapNewPlayerGrade map[uint8]uint64
}

// 鱼配置
type fishConfig struct {
	fishID uint32
	fishType uint32
}

type ConfigManager struct {
	paths map[uint32]*pathConfig
	drawRedPacket map[uint32]*roomDrawRedPacket
	fishConf map[uint32]uint32	// 鱼配置
}

func (mgr *ConfigManager) Load() {
	mgr.loadPathConfig()
	mgr.loadDrawRedPacketConfig()
	mgr.loadFishConfig()
	log.Println("fish's configs are loaded")
}

func (mgr *ConfigManager) readFile(path string) []byte {
	f, err := os.Open(path)
	util.CheckError(err)
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	util.CheckError(err)
	return data
}

func (mgr *ConfigManager) loadPathConfig()  {
	mgr.paths = make(map[uint32]*pathConfig)
	data := mgr.readFile("configs/fish/paths.json")
	jsv := make(map[string]interface{})
	err := json.Unmarshal(data, &jsv)
	util.CheckError(err)
	jvData, isOK := jsv["config"].([]interface{})
	if !isOK {
		util.CheckError(errors.New("paths's config is not existed or is not array"))
	}
	for _, v := range jvData {
		p, isOK := v.(map[string]interface{})
		if !isOK {
			util.CheckError(errors.New("paths's child of config is not object"))
		}
		conf := &pathConfig{}
		conf.ID = uint32(p["id"].(float64))
		conf.SwimDuration = uint32(p["swim_duration"].(float64))
		conf.StayDuration = uint32(p["stay_duration"].(float64))
		mgr.paths[conf.ID] = conf
	}
}

func (mgr *ConfigManager) getPathByID(id uint32) *pathConfig {
	if v, isOK := mgr.paths[id]; isOK {
		return v
	}
	return nil
}

func (mgr *ConfigManager) loadDrawRedPacketConfig() {
	mgr.drawRedPacket = make(map[uint32]*roomDrawRedPacket)
	data := mgr.readFile("configs/fish/red_packet_draw.json")
	jsv := make(map[string]interface{})
	err := json.Unmarshal(data, &jsv)
	util.CheckError(err)
	jvData, isOK := jsv["config"].([]interface{})
	if !isOK {
		util.CheckError(errors.New("red_packet_draw's config is not existed or is not array"))
	}
	for _, v := range jvData {
		p, isOK := v.(map[string]interface{})
		if !isOK {
			util.CheckError(errors.New("red_packet_draw's child of config is not object"))
		}
		conf := &roomDrawRedPacket{}
		conf.townID = uint32(p["town_id"].(float64))

		f := func(arr []interface{}, m map[uint8]uint64) {
			var (
				grade uint8
				chip uint64
			)
			for _, nv := range arr {
				np := nv.(map[string]interface{})
				grade = uint8(np["type"].(float64))
				chip = uint64(np["chipin"].(float64))
				m[grade] = chip
			}
		}
		normal := p["normal"].([]interface{})
		conf.mapNormalGrade = make(map[uint8]uint64)
		f(normal, conf.mapNormalGrade)
		newPlayer := p["new_player"].([]interface{})
		conf.mapNewPlayerGrade = make(map[uint8]uint64)
		f(newPlayer, conf.mapNewPlayerGrade)
		mgr.drawRedPacket[conf.townID] = conf
	}
}

func (mgr *ConfigManager) getTownDrawRedPacketByID(id uint32) *roomDrawRedPacket {
	if v, isOK := mgr.drawRedPacket[id]; isOK {
		return v
	}
	return nil
}

func (mgr *ConfigManager) loadFishConfig()  {
	mgr.fishConf = make(map[uint32]uint32)
	// 先读fish_res.json
	data := mgr.readFile("configs/fish/fish_res.json")
	//fishRes = make(map[uint32]uint32)
	jsv := make(map[string]interface{})
	err := json.Unmarshal(data, &jsv)
	util.CheckError(err)
	jvData, isOK := jsv["config"].([]interface{})
	if !isOK {
		util.CheckError(errors.New("fish_res's config is not existed or is not array"))
	}
	for _, v := range jvData {
		p, isOK := v.(map[string]interface{})
		if !isOK {
			util.CheckError(errors.New("fish_res's child of config is not object"))
		}
		kindID := uint32(p["kind_id"].(float64))
		fishType := uint32(p["type"].(float64))
		//fishRes[kindID] = fishType
		mgr.fishConf[kindID] = fishType

	}
	//// 再读room_fish.json
	//data = mgr.readFile("configs/fish/room_fish.json")
	//jsv = make(map[string]db_types{})
	//err = json.Unmarshal(data, &jsv)
	//util.CheckError(err)
	//jvData, isOK = jsv["config"].([]db_types{})
	//if !isOK {
	//	util.CheckError(errors.New("fish_res's config is not existed or is not array"))
	//}
	//for _, v := range jvData {
	//	p, isOK := v.(map[string]db_types{})
	//	if !isOK {
	//		util.CheckError(errors.New("fish_res's child of config is not object"))
	//	}
	//	fishID := uint32(p["id"].(float64))
	//	kindID := uint32(p["kind_id"].(float64))
	//	if t, isOK := fishRes[kindID]; isOK {
	//		mgr.fishConf[fishID] = t
	//	}
	//}
}

func (mgr *ConfigManager) getFishTypeByID(id uint32) uint32 {
	if v, isOK := mgr.fishConf[id]; isOK {
		return v
	}
	return 0
}