package mymysql

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github/go-robot/core/mydb/dbtypes"
	"log"
	"time"
)

type MySQL struct {
	src string
	sql *sql.DB
	ctx context.Context
}

func (my *MySQL) Open(ctx context.Context, acc string, pwd string, dbName string, addr string, port uint) {
	var err error
	my.src = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4,utf8", acc, pwd, addr, port, dbName)
	my.sql, err = sql.Open("mysql", my.src)
	if err != nil {
		log.Fatal("unable to use data source name", err)
		return
	}

	my.sql.SetConnMaxLifetime(time.Minute * 3)
	my.sql.SetMaxIdleConns(3)
	my.sql.SetMaxOpenConns(3)

	var stop func()
	my.ctx, stop = context.WithCancel(ctx)

	go func() {
		select {
		case <-ctx.Done():
			stop()
			my.sql.Close()
		}
	}()
}

func (my *MySQL)ping() bool {
	ctx, cancel := context.WithTimeout(my.ctx, 5 * time.Second)
	defer cancel()
	if err := my.sql.PingContext(ctx); err != nil {
		log.Panicln("ping error: ", err)
		return false
	}
	return true
}

func (my *MySQL)Query(sqlFmt string, args ...interface{}) []dbtypes.DBRow {
	if !my.ping() {
		return nil
	}
	ctx, cancel := context.WithTimeout(my.ctx, 5 * time.Second)
	defer cancel()
	stmt, err := my.sql.PrepareContext(ctx, sqlFmt)
	if err != nil {
		log.Panicln("query error: ", err)
		return nil
	}
	defer stmt.Close()
	ctx1, cancel1 := context.WithTimeout(my.ctx, 5 * time.Second)
	defer cancel1()
	rows, err := stmt.QueryContext(ctx1, args...)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Panicln("query error: ", err)
		}
		return nil
	}
	defer rows.Close()

	// 获取列
	cols, err := rows.Columns()
	if err != nil {
		log.Panicln("query get columns error: ", err)
		return nil
	}
	// 返回值引用列表
	scanArgs := make([]interface{}, len(cols))
	for i := range cols {
		var v interface{}
		scanArgs[i] = &v
	}
	// 结果列表
	results := make([]dbtypes.DBRow, 0)
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			log.Panicln("query result error: ", err)
			return nil
		}
		re := make(dbtypes.DBRow)
		for i, v := range scanArgs{
			re[cols[i]] = *v.(*interface{})
		}
		results = append(results, re)
	}
	return results
}

func (my *MySQL)Execute(sqlFmt string, args ...interface{}) (affectRows int64, lastInsertID int64) {
	if !my.ping() {
		return 0, 0
	}
	ctx, cancel := context.WithTimeout(my.ctx, 5 * time.Second)
	defer cancel()
	stmt, err := my.sql.PrepareContext(ctx, sqlFmt)
	if err != nil {
		log.Panicln("query error: ", err)
		return 0, 0
	}
	defer stmt.Close()
	ctx1, cancel1 := context.WithTimeout(my.ctx, 5 * time.Second)
	defer cancel1()
	re, err := stmt.ExecContext(ctx1, args...)
	if err != nil {
		log.Panicln("execute error: ", err)
		return 0, 0
	}
	affectRows, err = re.RowsAffected()
	if err != nil {
		log.Panicln("affect rows error:", err)
	}
	lastInsertID, err = re.LastInsertId()
	if err != nil {
		log.Panicln("last insert id error:", err)
	}
	return affectRows, lastInsertID
}
