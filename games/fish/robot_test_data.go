package fish

import (
	"context"
	"errors"
	"fmt"
	"github/go-robot/global"
	"github/go-robot/mydb"
	"log"
	"math/rand"
	"time"
)

func SetRobotTestData(ctx context.Context, pID uint32) bool {
	// db user
	userDB := getDBObject(ctx, global.GameCommonSetting.UserDB)
	if userDB == nil {
		log.Fatalln("fish test data need user db setting")
		return false
	}
	// db data
	dataDB := getDBObject(ctx, global.GameCommonSetting.DataDB)
	if dataDB == nil {
		log.Fatalln("fish test data need data db setting")
		return false
	}

	// 获取角色ID
	res := userDB.Query("select char_id from users where platform_id = ? and gz_id = ?", pID, global.MainSetting.GameZone)
	if len(res) != 1 {
		return false
	}
	charID := res[0]["char_id"].(int64)
	if charID == 0 {
		// 创建角色
		var err error
		charID, err = getCharID(userDB, dataDB, pID)
		if err != nil {
			return false
		}
	}
	// 设置物品
	itemSetting := global.FishTestDataSetting.Items
	if len(itemSetting) == 0 {
		return false
	}
	for id, num := range global.FishTestDataSetting.Items {
		it := dataDB.Query("select serial from game_items where container=? and model_id=?", charID, id)
		if len(it) == 0 {
			// insert
			afRows, _ := dataDB.Execute("insert into game_items set status=1, container=?, container_type=1, group_id=1, model_id=?, num=?",
				charID, id, num)
			if afRows == 0 {
				log.Printf("set item failed: mode_id=%d, num=%d", id, num)
				return false
			}
			continue
		}
		dataDB.Execute("update game_items set num=? where serial=?", num, it[0]["serial"].(int64))
	}
	// 炮等级
	dataDB.Execute("update game_chars set official=? where serial=?", global.FishTestDataSetting.CaliberLV, charID)
	return true
}

func getDBObject(ctx context.Context, dbSetting global.DBSetting) mydb.MyDB {
	if !dbSetting.IsUsable {
		return nil
	}
	db := mydb.NewDB(ctx, global.MySQL, dbSetting.Account, dbSetting.Password, dbSetting.Database, dbSetting.Address, dbSetting.Port)
	if db == nil {
		return nil
	}
	return db
}

func getCharID(userDB mydb.MyDB, dataDB mydb.MyDB, pID uint32) (charID int64, err error) {
	// 获取角色ID
	res := userDB.Query("select uid, nick, char_id from users where platform_id = ? and gz_id = ?", pID, global.MainSetting.GameZone)
	if len(res) != 1 {
		return 0, errors.New("don't find user in users table")
	}
	row := res[0]
	charID = row["char_id"].(int64)
	if charID == 0 {
		// 创建角色
		charID, err = createCharacter(userDB, dataDB, row["uid"].(int64), row.GetString("nick"))
	}
	return charID, err
}

func createCharacter(userDB mydb.MyDB, dataDB mydb.MyDB, uID int64, nick string) (int64, error) {
	x := rand.Intn(791) + 100
	y := rand.Intn(791) + 100
	id := 0
	if rand.Intn(2) == 0 {
		id = rand.Intn(3) + 511
	} else {
		id = rand.Intn(3) + 411
	}
	afRows, lastInsertID := dataDB.Execute("insert into game_chars set x=?, y=?, lv=1, id=?, name=?, userid=?, " +
		"status=?, create_time=?", x, y, id, nick, uID, 1, time.Now().Unix())
	if afRows == 0 {
		return 0, errors.New(fmt.Sprintf("create char error: userid = %d", uID))
	}
	// update users
	userDB.Execute("update users set char_id = ? where uid = ?", lastInsertID, uID)
	return lastInsertID, nil
}