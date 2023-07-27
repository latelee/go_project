/*
go test -v -run TestSqlite

输出信息 无数据库文件
test of sqlte3...
connect to  foobar.db3 ok
run sql:
select version, updateTime from myVersion order by version desc limit 1
not found table, will create it.
got db version [] update time []
connect to  foobar.db3 ok
insert db version [] at: [2023-12-02 10:42:18]
insert result:  <nil>
--- PASS: TestSqlite (1.04s)
PASS

已有数据但版本较新
test of sqlte3...
connect to  foobar.db3 ok
run sql: [select version, updateTime from myVersion order by version desc limit 1]
got db version [20231202] update time [2023-12-02T10:48:20Z]
connect to  foobar.db3 ok
insert db version [20231203] at: [2023-12-02 10:48:47]
insert result:  <nil>
--- PASS: TestSqlite (1.03s)
PASS
*/

package test

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
	"webdemo/pkg/com"

	_ "github.com/mattn/go-sqlite3"
)

var (
	// 数据库文件名及表名
	dbServer     string = "foobar.db3"
	tableVersion string = "myVersion"
	tableList    string = "myList"
)

// 信息表 结构体可对于json风格数据传输解析
type InfoList_t struct {
	Id         int    `json:"-"`
	Version    string `json:"-"`
	Name       string `json:"-"`
	City       string `json:"-"`
	UpdateTime string `json:"-"`
}

var sqlarr []string = []string{
	// 版本号
	`CREATE TABLE "myVersion" (
		"version" VARCHAR(20) NOT NULL,
		"updateTime" datetime DEFAULT "",
		PRIMARY KEY ("version")
	);`,

	// 信息表
	`CREATE TABLE "myList" (
		"id" int NOT NULL,
		"version" VARCHAR(20) NOT NULL,
		"name" VARCHAR(20) NOT NULL,
		"city" VARCHAR(20) NOT NULL,
		"updateTime" datetime DEFAULT "",
		PRIMARY KEY ("id")
	);`,
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func Sleep(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func CreateSqlite3(dbname string, create bool) (sqldb *sql.DB, err error) {
	if create == false && !IsExist(dbname) {
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
	fmt.Println("connect to ", dbname, "ok")

	return
}

func readOrCreateDBTable(sqldb *sql.DB) (version, updateTime string, err error) {
	needCreate := false
	sqlstr := fmt.Sprintf(`select version, updateTime from %v order by version desc limit 1`,
		tableVersion)
	fmt.Printf("run sql: [%v]\n", sqlstr)
	var results *sql.Rows
	results, err = sqldb.Query(sqlstr)
	if err != nil {
		if strings.Contains(err.Error(), "no such table") {
			needCreate = true
		} else {
			fmt.Println("query error: ", err)
			return
		}
	}

	if !needCreate {
		for results.Next() {
			var item1, item2 sql.NullString
			err = results.Scan(&item1, &item2)
			if err != nil {
				fmt.Println("scan error: ", err)
				break
			}
			if !item1.Valid || !item2.Valid {
				continue
			}
			version = item1.String
			updateTime = item2.String
		}

		defer results.Close()
	} else {
		fmt.Println("not found table, will create it.")
		for _, item := range sqlarr {
			_, err = sqldb.Exec(item)
			if err != nil {
				fmt.Printf("Exec sql failed: [%v] [%v] \n", err, item)
			}
		}
	}

	// return
	columnName := "age" // mileage
	sqlstr = fmt.Sprintf(`select sql from sqlite_master where type = 'table' and name = '%v' and sql like '%%%v%%'`,
		tableList, columnName)
	fmt.Printf("check column run sql: [%v]\n", sqlstr)
	results, err = sqldb.Query(sqlstr)
	if err != nil {
		fmt.Println("query error: ", err)
		return
	}
	if results.Next() == false {
		fmt.Printf("not found %v, will add it\n", columnName)
		sqlstr = fmt.Sprintf(`ALTER TABLE %v ADD COLUMN %v INTEGER DEFAULT -1`, tableList, columnName)
		_, err = sqldb.Exec(sqlstr)
		if err != nil {
			fmt.Printf("Exec sql failed: [%v] [%v] \n", err, sqlstr)
		}
	}

	return
}

func readTableVersion(sqldb *sql.DB) (version, updateTime string, err error) {
	sqlstr := fmt.Sprintf(`select version, updateTime from %v order by version desc limit 1`,
		tableVersion)
	fmt.Printf("run sql: [%v]\n", sqlstr)
	var results *sql.Rows
	results, err = sqldb.Query(sqlstr)
	if err != nil {
		fmt.Println("query error: ", err)
		return
	}
	defer results.Close()

	for results.Next() {
		var item1, item2 sql.NullString
		err = results.Scan(&item1, &item2)
		if err != nil {
			fmt.Println("scan error: ", err)
			break
		}
		if !item1.Valid || !item2.Valid {
			continue
		}
		version = item1.String
		updateTime = item2.String
	}

	return
}

func insertDBDetail(tx *sql.Tx, gxList []InfoList_t, version string) (err error) {
	tablename := tableList
	sqlstr := fmt.Sprintf(`DELETE FROM %v`, tablename)
	stmt, err := tx.Prepare(sqlstr)
	if err != nil {
		err = errors.New("prepare for [" + sqlstr + "] failed: " + err.Error())
		return
	}
	_, err = stmt.Exec()
	if err != nil {
		err = errors.New("delete " + tablename + "failed: " + err.Error())
		return
	}

	sqlstr = fmt.Sprintf(`INSERT OR REPLACE INTO %v 
(id, version, name, city, updateTime) 
VALUES (?, ?, ?, ?, ?)`,
		tablename)
	stmt, _ = tx.Prepare(sqlstr)
	for _, item := range gxList {
		// item.Id = idx
		item.Version = version
		item.UpdateTime = com.GetNowDateTime("YYYY-MM-DD HH:mm:ss")
		_, err = stmt.Exec(item.Id, item.Version, item.Name, item.City, item.UpdateTime)
		if err != nil {
			err = errors.New("insert " + tablename + "failed: " + err.Error())
			return
		}
	}

	return
	// debug 制作bug
	// TODO 制作锁住，制作语法错误
	err = errors.New("database is locked")

	return
}

func insertDBVersion(tx *sql.Tx, version string) (err error) {
	tablename := tableVersion
	sqlstr := fmt.Sprintf(`DELETE FROM %v`, tablename)
	stmt, err := tx.Prepare(sqlstr)
	if err != nil {
		err = errors.New("prepare for [" + sqlstr + "] failed: " + err.Error())
		return
	}
	_, err = stmt.Exec()
	if err != nil {
		err = errors.New("delete " + tablename + " failed: " + err.Error())
		return
	}

	sqlstr = fmt.Sprintf(`INSERT OR REPLACE INTO %v (version, updateTime) VALUES (?, ?)`, tablename)
	stmt, err = tx.Prepare(sqlstr)
	if err != nil {
		err = errors.New("prepare for [" + sqlstr + "] failed: " + err.Error())
		return
	}
	updateTime := com.GetNowDateTime("YYYY-MM-DD HH:mm:ss")
	fmt.Printf("insert db version [%v] at: [%v]\n", version, updateTime)
	_, err = stmt.Exec(version, updateTime)
	if err != nil {
		err = errors.New("insert " + tablename + "failed: " + err.Error())
		return
	}

	return
}

// 入库2个表，以事务方式
func insertDBBatch(gxList []InfoList_t, version string) (err error) {
	SQLDB, err := CreateSqlite3(dbServer, false)
	if err != nil {
		// fmt.Println(err.Error())
		return err
	}

	var tx *sql.Tx
	tx, err = SQLDB.Begin()
	if err != nil {
		err = errors.New("begin sql error: " + err.Error())
		return err
	}

	defer func() {
		if err != nil {
			err = errors.New("exec sql failed rollback: " + err.Error())
			tx.Rollback()
		} else {
			err = nil
			tx.Commit()
		}
		// 延时一会，关闭
		Sleep(1000)
		SQLDB.Close()
	}()

	err = insertDBVersion(tx, version)
	if err != nil {
		return
	}

	err = insertDBDetail(tx, gxList, version)
	if err != nil {
		return
	}

	return
}

//////////////////////
func makeData() (gxList []InfoList_t) {
	var tmp InfoList_t
	tmp.Id = 100
	tmp.Version = "100"
	tmp.Name = "latelee"
	tmp.City = "梧州"
	gxList = append(gxList, tmp)

	tmp = InfoList_t{}
	tmp.Id = 250
	tmp.Version = "250"
	tmp.Name = "latelee"
	tmp.City = "岑溪"
	gxList = append(gxList, tmp)

	return
}

// 读取基础信息，尝试创建表
func readDBVersion(atype int) (version, datetime string) {
	SQLDB, err := CreateSqlite3(dbServer, true)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if atype == 0 {
		version, datetime, _ = readOrCreateDBTable(SQLDB)
	} else {
		version, datetime, _ = readTableVersion(SQLDB)
	}
	SQLDB.Close()

	return
}
func TestSqlite(t *testing.T) {
	fmt.Println("test of sqlte3...")

	// 1 尝试获取数据表的版本号（可能为空）
	version, datetime := readDBVersion(0)
	fmt.Printf("init: got db version [%v] update time [%v]\n", version, datetime)

	// return
	cnt := 0
	myVer := 10000
	if version != "" {
		myVer, _ = strconv.Atoi(version)
	}

	for {
		// 2 模拟业务：自定义版本号，较新时，才入库
		version, datetime := readDBVersion(1)
		fmt.Printf("got db version [%v] update time [%v]\n", version, datetime)
		myVer += 1
		newVer := fmt.Sprintf("%v", myVer)
		if newVer > version {
			data := makeData()
			err := insertDBBatch(data, newVer)
			fmt.Println("insert result: ", err)
		} else {
			fmt.Println("db is newest, do nothing")
		}
		cnt++
		if cnt > 1 {
			break
		}
		Sleep(1000)
	}
	fmt.Println("test done.")
}
