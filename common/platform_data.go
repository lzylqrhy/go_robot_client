package common

import (
	"encoding/json"
	"fmt"
	"github/go-robot/util"
	"io/ioutil"
	"net/http"
)

type PlatformData struct {
	PID        uint32
	Nickname   string
	LoginToken string
	ServerAddr string
}

func GetPlatformUserData(start uint, num uint) []*PlatformData {
	end := start + num - 1	// 如果相等，则取当前
	sUrl := fmt.Sprintf("http://app_fish.dev.com/platform/genRegisteredGameRobot?start=%d&end=%d&vaild=1",start,end)
	resp, err := http.Get(sUrl)
	util.CheckError(err)
	body, err := ioutil.ReadAll(resp.Body)
	util.CheckError(err)
	jv := make(map[string]interface{})
	err = json.Unmarshal(body, &jv)
	util.CheckError(err)

	userList := make([]*PlatformData, 0, end-start)
	ret := jv["ret"].(float64)
	if 1 == uint32(ret) {
		jvData := jv["data"].([]interface{})
		for _, v := range jvData {
			p := v.(map[string]interface{})
			pd := new(PlatformData)
			pd.PID = uint32(p["uid"].(float64))
			pd.Nickname = p["nickname"].(string)
			pd.LoginToken = p["loginkey"].(map[string]interface{})["783"].(string)
			pd.ServerAddr = p["ws_host"].(map[string]interface{})["783"].(string)
			userList = append(userList, pd)
		}
	}
	return userList
}