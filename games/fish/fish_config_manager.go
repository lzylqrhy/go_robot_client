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

type ConfigManager struct {
	paths map[uint32]*pathConfig
}

func (mgr *ConfigManager) Load() {
	mgr.loadPathConfig()
	log.Println("fish's configs are loaded")
}

func (mgr *ConfigManager) readFile(path string) []byte {
	f, err := os.Open("./configs/fish/paths.json")
	util.CheckError(err)
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	util.CheckError(err)
	return data
}

func (mgr *ConfigManager) loadPathConfig()  {
	mgr.paths = make(map[uint32]*pathConfig)
	data := mgr.readFile("configs/paths.json")
	jsv := make(map[string]interface{})
	err := json.Unmarshal(data, &jsv)
	util.CheckError(err)
	jvData, isOK := jsv["config"].([]interface{})
	if !isOK {
		util.CheckError(errors.New("config is not existed or is not array"))
	}
	for _, v := range jvData {
		p, isOK := v.(map[string]interface{})
		if !isOK {
			util.CheckError(errors.New("child of config is not object"))
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