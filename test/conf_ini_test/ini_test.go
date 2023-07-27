/*
go test -v -run TestIni

输出：
test of ini...
keystrings: [sec1 sec2]
key: [sec1]
value: [110 210 ddd 99] [133 135 1 2 1588 1509] [310-410]
--------
key: [sec2]
value: [late \n lee] [Jim Kent]
--------
section2.id value: 250
section2.arrint value: [1 2 3]
section2.arrstr value: [apple apple2 apple3]

*/

package test

import (
	"fmt"
	"testing"

	"gopkg.in/ini.v1"
)

var (
	cfgFile string = "config.ini"
)

func TestIni(t *testing.T) {
	fmt.Println("test of ini...")
	// 加载
	cfg, err := ini.Load(cfgFile)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		return
	}

	// 读一节
	mySection := cfg.Section("section1")
	fmt.Println("keystrings:", mySection.KeyStrings())
	// 解析
	for _, key := range mySection.KeyStrings() {

		fmt.Printf("key: [%v]\nvalue: ", key)
		values := mySection.Key(key).Strings("|")
		for _, item := range values {
			fmt.Printf("[%v] ", item)
		}
		fmt.Printf("\n--------\n")
	}

	// 按单个字符串
	strVlaue := cfg.Section("section2").Key("id").Value()
	fmt.Printf("section2.id value: %v\n", strVlaue)

	// 数值，以"|"分隔
	arrVlaueInt := cfg.Section("section2").Key("arrint").ValidInts("|")
	fmt.Printf("section2.arrint value: %v\n", arrVlaueInt)

	// 字符串，以","分隔
	arrVlaueStr := cfg.Section("section2").Key("arrstr").Strings(",")
	fmt.Printf("section2.arrstr value: %v\n", arrVlaueStr)

	// 写入ini
	cfg.Section("").Key("name").SetValue("foobar")
	cfg.Section("info").Key("appVer").SetValue("1.3")
	cfg.Section("info").Key("author").SetValue("Late Lee")

	// 另起文件
	cfg.SaveTo("output.ini")
}
