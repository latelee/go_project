package db

import (
	//"database/sql"
	//"flag"
	// "fmt"
    "log"
    //"errors"
    "time"
    //"reflect"
    //"math"
    //"strconv"
    // 导入mysql驱动
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/denisenkom/go-mssqldb"
    _ "github.com/mattn/go-oci8"
    _ "github.com/mattn/go-sqlite3"
    "github.com/go-xorm/xorm"
    "github.com/go-xorm/core"
    //"strings"
    //"encoding/binary"
	//"io/ioutil"
	
	// conf "github.com/latelee/dbtool/common/conf"
	
)

type DBParam struct {
    windows     bool    // 是否windows身份认证
	server      string
	database    string
    user        string
	passwd      string
}

///////////////

func CreateSqlServerXorm(dbstr string) (engine *xorm.Engine) {
    engine, err := CreateXorm("mssql", dbstr)
	if err != nil {
		log.Println("Open database failed:", err)
		return nil
	}

    return
}

func CreateMysqlXorm(dbstr string) (engine *xorm.Engine) {
    engine, err := CreateXorm("mysql", dbstr)
	if err != nil {
		log.Println("Open database failed:", err)
		return nil
	}
    //log.Println("connect to ", dbParam.server, dbParam.database, "ok")
    
    return
}

func CreateSqlite3Xorm(dbname string) (engine *xorm.Engine) {
    engine, err := CreateXorm("sqlite3", dbname)
	if err != nil {
		log.Println("Open database failed:", err)
		return nil
	}
    log.Println("connect to ", dbname, "ok")

    return
}

func CreateOracleXorm(dbstr string) (engine *xorm.Engine) {
    engine, err := CreateXorm("oci8", dbstr)
	if err != nil {
		log.Println("Open database failed:", err)
		return nil
	}
    //log.Println("connect to ", dbParam.server, dbParam.database, "ok")
	
    
    return
}

func CreateXorm(dbType, connString string) (engine *xorm.Engine, err error) {

	engine, err = xorm.NewEngine(dbType, connString)
	if err != nil {
		return
	}

	// 注：这里设置为UTC时间，待确认具体的
	engine.DatabaseTZ = time.UTC  // time.Local
	engine.TZLocation = time.UTC

	engine.SetMaxIdleConns(4)
	engine.SetMaxOpenConns(32)

    // set same name...
    engine.SetMapper(core.SameMapper{})
    
	//engine.ShowSQL(true)
    engine.ShowSQL(false)
    
    return
}