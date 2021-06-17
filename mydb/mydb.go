package mydb

import (
	"context"
	"github/go-robot/global"
	"github/go-robot/mydb/dbtypes"
	"github/go-robot/mydb/mymysql"
)

type MyDB interface {
	Query(sqlFmt string, args ...interface{}) [] dbtypes.DBRow
	Execute(sqlFmt string, args ...interface{}) (affectRows int64, lastInsertID int64)
}

func NewDB(ctx context.Context, driver uint, acc string, pwd string, dbName string, addr string, port uint) MyDB {
	switch driver {
	case global.MySQL:
		db := new(mymysql.MySQL)
		db.Open(ctx, acc, pwd, dbName, addr, port)
		return db
	}
	return nil
}

