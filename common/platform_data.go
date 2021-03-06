/**
 这个文件的用于处理和平台交互的公共逻辑
 created by lzy
*/
package common

import (
	"encoding/json"
	"fmt"
	"github/go-robot/global/ini"
	"github/go-robot/util"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// 机器人平台数据类型
type PlatformData struct {
	PID          uint32
	Nickname     string
	LoginToken   string
	WSServerAddr string
}
// 获取机器人信息列表
func GetPlatformUserData() []*PlatformData {
	cfg := &ini.MainSetting
	if 0 == cfg.RobotNum {
		return nil
	}
	var sUrl string
	if strings.Index(cfg.RobotAddr, "?") != -1 {
		sUrl = fmt.Sprintf("%s&start=%d&end=%d&vaild=1",
			cfg.RobotAddr, cfg.RobotStart, cfg.RobotStart + cfg.RobotNum - 1)
	} else {
		sUrl = fmt.Sprintf("%s?start=%d&end=%d&vaild=1",
			cfg.RobotAddr, cfg.RobotStart, cfg.RobotStart + cfg.RobotNum - 1)
	}
	resp, err := http.Get(sUrl)
	util.CheckError(err)
	body, err := ioutil.ReadAll(resp.Body)
	util.CheckError(err)
	jv := make(map[string]interface{})
	err = json.Unmarshal(body, &jv)
	util.CheckError(err)
	// 用户列表
	userList := make([]*PlatformData, 0, cfg.RobotNum)
	ret := jv["ret"].(float64)
	if 1 == uint32(ret) {
		jvData := jv["data"].([]interface{})
		for _, v := range jvData {
			p := v.(map[string]interface{})
			pd := new(PlatformData)
			pd.PID = uint32(p["uid"].(float64))
			pd.Nickname = p["nickname"].(string)
			pd.LoginToken = p["loginkey"].(map[string]interface{})[cfg.GameZone].(string)
			pd.WSServerAddr = p["ws_host"].(map[string]interface{})[cfg.GameZone].(string)
			userList = append(userList, pd)
		}
	} else {
		log.Println("access url failed, data:", string(body))
	}
	return userList
}