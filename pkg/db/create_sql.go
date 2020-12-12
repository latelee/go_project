package db

import (
	"database/sql"
	//"os"
	// "fmt"
	"log"
	"errors"
    //"errors"
    //"time"
    //"reflect"
    //"math"
    //"strconv"
    // 导入mysql驱动
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/denisenkom/go-mssqldb"
    _ "github.com/mattn/go-oci8"
    _ "github.com/mattn/go-sqlite3"
    //"github.com/go-xorm/xorm"
    //"github.com/go-xorm/core"
    //"strings"
    //"encoding/binary"
	//"io/ioutil"

	"github.com/latelee/dbtool/pkg/com"
	// conf "github.com/latelee/dbtool/common/conf"
)

func CreateSqlServer(dbstr string) (sqldb *sql.DB, err error) {
    sqldb, err = sql.Open("mssql", dbstr)
	if err != nil {
		return nil, errors.New("open database failed: " + err.Error())
	}
	// Open不一定会连接数据库，Ping可能会连接
	err = sqldb.Ping()
	if err != nil {
		return nil, errors.New("connect database failed: " + err.Error())
	}
	log.Println("connect to sqlserver ok")
    //log.Println("connect to ", dbParam.server, dbParam.database, "ok")

    return sqldb, nil
}

func CreateMysql(dbstr string) (sqldb *sql.DB, err error) {
    sqldb, err = sql.Open("mysql", dbstr)
	if err != nil {
		return nil, errors.New("open database failed: " + err.Error())
	}
	err = sqldb.Ping()
	if err != nil {
		return nil, errors.New("connect database failed: " + err.Error())
	}
    log.Println("connect to mysql ok")

    return
}

func CreateSqlite3(dbname string) (sqldb *sql.DB, err error) {
	if !com.IsExist(dbname) {
		return nil, errors.New("open database failed: " + dbname + " not found")
	}
    sqldb, err = sql.Open("sqlite3", dbname)
	if err != nil {
		return nil, errors.New("open database failed: " + err.Error())
	}
	err = sqldb.Ping()
	if err != nil {
		return nil, errors.New("connect database failed: " + err.Error())
	}
    log.Println("connect to ", dbname, "ok")

    return
}

func CreateOracle(dbstr string) (sqldb *sql.DB, err error) {
    sqldb, err = sql.Open("oci8", dbstr)
	if err != nil {
		return nil, errors.New("open database failed: " + err.Error())
	}
	err = sqldb.Ping()
	if err != nil {
		return nil, errors.New("connect database failed: " + err.Error())
	}
	log.Println("connect to oracle ok")
    //log.Println("connect to ", dbParam.server, dbParam.database, "ok")
    
    return
}
