package fish

var ConfMgr = ConfigManager{}

// 在main之前执行
func init() {
	ConfMgr.Load()
}

const BigFish = 3