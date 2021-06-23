package global

// 日志类型
const (
	Debug = 1 << iota
	Warning
	Error
	Fatal
	AllLog = Debug|Warning|Error|Fatal
)

// 游戏类型
const (
	FishGame = 1 + iota
	FruitGame
	AladdinGame
)

// 物品
const (
	ItemCoin = 1000	// 金币
	ItemFishBlackMissile = 2001	//  青铜炸弹
	ItemFishBronzeMissile = 2002	// 黑铁炸弹
	ItemFishSilverMissile = 2003	// 白银炸弹
	ItemFishGoldMissile = 2004	// 黄金炸弹
	ItemFishPlatinumMissile = 2005	// 铂金炸弹
	ItemFishKingMissile = 2006	// 至尊炸弹
)

// 数据库
const (
	MySQL = 1 + iota
)
