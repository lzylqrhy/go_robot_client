package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	// 读取配置
	servAddr := "192.168.0.194:7710"
	// 机器人数量
	playerNum := 1
	// 从平台获取信息
	httpRequest()
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	for i := 0; i < playerNum; i++ {
		wg.Add(1)
		subCtx := context.WithValue(ctx, "index", i)
		go createConnect(subCtx, &wg, servAddr)
	}
	// 监听信号
	waitForASignal()
	cancel()
	fmt.Println("stop all jobs")
	wg.Wait()
	fmt.Println("exit")
}

func waitForASignal()  {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)
	<-sig
}


type AutoGenerated struct {
	Ret       int    `json:"ret"`
	Status    int    `json:"status"`
	Msg       string `json:"msg"`
	ClientMsg string `json:"clientMsg"`
	Data      []struct {
		UID                   int    `json:"uid"`
		SerialID              int    `json:"serial_id"`
		Luid                  int    `json:"luid"`
		ClientUID             string `json:"client_uid"`
		Agentid               int    `json:"agentid"`
		Username              string `json:"username"`
		Gender                int    `json:"gender"`
		Nickname              string `json:"nickname"`
		UserNickname          string `json:"user_nickname"`
		Password              string `json:"password"`
		Openid                string `json:"openid"`
		Openkey               string `json:"openkey"`
		Unionid               string `json:"unionid"`
		WxServerOpenid        string `json:"wx_server_openid"`
		MainPassword          string `json:"main_password"`
		ThumbID               int    `json:"thumb_id"`
		ThumbURL              string `json:"thumb_url"`
		UserIconURL           string `json:"user_icon_url"`
		BindPhone             string `json:"bind_phone"`
		BindPhoneTime         int    `json:"bind_phone_time"`
		BindEmail             string `json:"bind_email"`
		BindEmailTime         int    `json:"bind_email_time"`
		VerifyCode            string `json:"verify_code"`
		RegTime               int    `json:"reg_time"`
		BanTime               int    `json:"ban_time"`
		BanMsg                string `json:"ban_msg"`
		RegResourceID         int    `json:"reg_resource_id"`
		RegIP                 string `json:"reg_ip"`
		FakeIP                string `json:"fake_ip"`
		RegDeviceID           int    `json:"reg_device_id"`
		ActiveChannelDeviceID int    `json:"active_channel_device_id"`
		LastDeviceID          int    `json:"last_device_id"`
		RegChannel            int    `json:"reg_channel"`
		ActiveChannel         int    `json:"active_channel"`
		ActiveChannelStatus   int    `json:"active_channel_status"`
		LastIP                string `json:"last_ip"`
		SafeEmail             string `json:"safe_email"`
		SafePhone             string `json:"safe_phone"`
		LoginNum              int    `json:"login_num"`
		LoginTime             int    `json:"login_time"`
		Status                int    `json:"status"`
		MobileUser            int    `json:"mobile_user"`
		WebUser               int    `json:"web_user"`
		MobileLoginTime       int    `json:"mobile_login_time"`
		LinkTime              int    `json:"link_time"`
		Pf                    int    `json:"pf"`
		Via                   int    `json:"via"`
		ViaType               int    `json:"via_type"`
		ParentUID             int    `json:"parent_uid"`
		ReferUID              int    `json:"refer_uid"`
		ReferTime             int    `json:"refer_time"`
		ShopID                int    `json:"shop_id"`
		ShopLevel             int    `json:"shop_level"`
		Info                  string `json:"info"`
		AccessKey             string `json:"access_key"`
		WxServiceOpenid       string `json:"wx_service_openid"`
		NewbieGift            int    `json:"newbie_gift"`
		TrailClient           string `json:"trail_client"`
		Modified              int    `json:"modified"`
		PayTotal              string `json:"pay_total"`
		AlipayAccount         string `json:"alipay_account"`
		AlipayRealName        string `json:"alipay_real_name"`
		Loginkey              struct {
			Num780 string `json:"780"`
			Num781 string `json:"781"`
			Num782 string `json:"782"`
			Num783 string `json:"783"`
		} `json:"loginkey"`
		WsHost struct {
			Num780 string `json:"780"`
			Num781 string `json:"781"`
			Num782 string `json:"782"`
			Num783 string `json:"783"`
		} `json:"ws_host"`
		WssHost struct {
			Num780 string `json:"780"`
			Num781 string `json:"781"`
			Num782 string `json:"782"`
			Num783 string `json:"783"`
		} `json:"wss_host"`
	} `json:"data"`
}

type PlatformData struct {
	PID int
	Nickname string
}

func httpRequest() {
	resp, err := http.Get("http://app_fish.dev.com/platform/genRegisteredGameRobot?start=2&end=2&vaild=1")
	checkError(err)
	body, err := ioutil.ReadAll(resp.Body)
	checkError(err)
	jv := make(map[string]interface{})
	err = json.Unmarshal(body, &jv)
	checkError(err)

	var pd PlatformData
	ret, bOk := jv["ret"].(int)
	if bOk && 1 == ret {
		jvData := jv["data"].([]interface{})
		for _, v := range jvData {
			p := v.(map[string]interface{})
			pd.PID = p["client_uid"].(string)
		}
	}
	str := string(body)
	fmt.Println(str)
}

func checkError(err error)  {
	if err != nil {
		panic(err)
	}
}