package fish

var ConfMgr = ConfigManager{}

func init() {
	ConfMgr.Load()
}