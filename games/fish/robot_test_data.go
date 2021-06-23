package fish

import (
	"context"
	"errors"
	"fmt"
	"github/go-robot/common"
	"github/go-robot/core/mydb"
	"github/go-robot/global"
	"github/go-robot/global/ini"
	"log"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// 多协程设置数据
func RunTestData(pds []*common.PlatformData) {
	uLen := len(pds)
	if uLen == 0 {
		return
	}
	num := runtime.NumCPU()
	if uLen < num {
		num = uLen
	}
	wg := sync.WaitGroup{}
	ctx := context.Background()
	index := int32(-1)
	atomic.LoadInt32(&index)
	for i := 0; i < num; i++ {
		wg.Add(1)
		op := operateData{}
		op.Init(ctx)
		go func(){
			defer wg.Done()
			for {
				i := atomic.AddInt32(&index, 1)
				if int(i) >= uLen {
					break
				}
				if !op.SetTestData(pds[i].PID) {
					log.Printf("set data of robot failed, index=%d pid=%d", i, pds[i].PID)
				} else {
					log.Printf("set data of robot successfully, index=%d pid=%d", i, pds[i].PID)
				}
			}
		}()
	}
	wg.Wait()
}

// 处理数据类
type operateData struct {
	userDB mydb.MyDB
	dataDB mydb.MyDB
}

func (o *operateData)Init(ctx context.Context) {
	// db user
	o.userDB = getDBObject(ctx, ini.GameCommonSetting.UserDB)
	if o.userDB == nil {
		log.Fatalln("fish test data need valid user db setting, please check config")
		return
	}
	// db data
	o.dataDB = getDBObject(ctx, ini.GameCommonSetting.DataDB)
	if o.dataDB == nil {
		log.Fatalln("fish test data need valid data db setting, please check config")
		return
	}
}

func (o *operateData) SetTestData(pID uint32) bool {
	// 获取角色ID
	res := o.userDB.Query("select char_id from users where platform_id = ? and gz_id = ?", pID, ini.MainSetting.GameZone)
	if len(res) != 1 {
		log.Printf("find char_id failed from users table : platform_id=%d, gz_id=%s", pID, ini.MainSetting.GameZone)
		return false
	}
	charID := res[0]["char_id"].(int64)
	if charID == 0 {
		// 创建角色
		var err error
		charID, err = o.getCharID(pID)
		if err != nil {
			return false
		}
	}
	// 设置物品
	itemSetting := ini.FishTestDataSetting.Items
	if len(itemSetting) == 0 {
		return false
	}
	for id, num := range ini.FishTestDataSetting.Items {
		it := o.dataDB.Query("select serial from game_items where container=? and model_id=?", charID, id)
		if len(it) == 0 {
			// insert
			afRows, _ := o.dataDB.Execute("insert into game_items set status=1, container=?, container_type=1, group_id=1, model_id=?, num=?",
				charID, id, num)
			if afRows == 0 {
				log.Printf("set item failed: mode_id=%d, num=%d", id, num)
				return false
			}
			continue
		}
		o.dataDB.Execute("update game_items set num=? where serial=?", num, it[0]["serial"].(int64))
	}
	// 炮等级
	o.dataDB.Execute("update game_chars set official=? where serial=?", ini.FishTestDataSetting.CaliberLV, charID)
	return true
}

func getDBObject(ctx context.Context, dbSetting ini.DBSetting) mydb.MyDB {
	if !dbSetting.IsUsable {
		return nil
	}
	db := mydb.NewDB(ctx, global.MySQL, dbSetting.Account, dbSetting.Password, dbSetting.Database, dbSetting.Address, dbSetting.Port)
	if db == nil {
		return nil
	}
	return db
}

func (o *operateData)getCharID(pID uint32) (charID int64, err error) {
	// 获取角色ID
	res := o.userDB.Query("select uid, nick, char_id from users where platform_id = ? and gz_id = ?", pID, ini.MainSetting.GameZone)
	if len(res) != 1 {
		return 0, errors.New("don't find user in users table")
	}
	row := res[0]
	charID = row["char_id"].(int64)
	if charID == 0 {
		// 创建角色
		charID, err = o.createCharacter(row["uid"].(int64), row.GetString("nick"))
	}
	return charID, err
}

func (o *operateData)createCharacter(uID int64, nick string) (int64, error) {
	x := rand.Intn(791) + 100
	y := rand.Intn(791) + 100
	id := 0
	if rand.Intn(2) == 0 {
		id = rand.Intn(3) + 511
	} else {
		id = rand.Intn(3) + 411
	}
	afRows, lastInsertID := o.dataDB.Execute("insert into game_chars set x=?, y=?, lv=1, id=?, name=?, userid=?, " +
		"status=?, create_time=?", x, y, id, nick, uID, 1, time.Now().Unix())
	if afRows == 0 {
		return 0, errors.New(fmt.Sprintf("create char error: userid = %d", uID))
	}
	// update users
	o.userDB.Execute("update users set char_id = ? where uid = ?", lastInsertID, uID)
	return lastInsertID, nil
}