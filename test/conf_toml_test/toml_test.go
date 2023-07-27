/*
go test -v -run TestToml

输出：
test of toml...
need: true name: conf file
version: 2.0
time: 2020-10-03 09:21:13 +0000 UTC
empty: empty
text:   hello
  world!

name: late
 lee, name1:  age: 99
sta[0]: jim kent jk@latelee.org
sta[1]: late lee li@latelee.org
sta[2]: foo foo@latelee.org
fruit: [apple apple1 apple2 apple3 apple4 apple5]
ports: [8080 8081 8082]
bad: [] bad1: [0]
logdir: log
name: 在线 url: http://abc.com
name: 离线 url: http://ccc.com

*/

package test

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"
)

var (
	cfgFile string // = "config.toml"
)

type mapUrl_t struct {
	Name  string `json:"name"`
	Attri string `json:"attri"`
	Url   string `json:"url"`
}

func TestToml(t *testing.T) {
	fmt.Println("test of toml...")

	// 设置配置文件的2种方式
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("./")
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
	}

	// 读取 注：如果toml格式有误，此处报错
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("'%v' file read error: %v\n", cfgFile, err)
		return
	}

	name := viper.GetString("name") // 读取 字符串
	version := viper.GetString("version")

	need := viper.GetBool("need") // 读取 布尔
	theTime := viper.GetString("time")
	empty := viper.GetString("empty")
	text := viper.GetString("text")

	fmt.Printf("need: %v name: %v\nversion: %v \ntime: %v \nempty: %s \ntext: %v\n", need, name, version, theTime, empty, text)

	// 多级读取
	name = viper.GetString("my.name")
	name1 := viper.GetString("my.name1")
	age := viper.GetInt("my.age")
	fmt.Printf("name: %v, name1: %v age: %v \n", name, name1, age)

	// 字符串数组
	newSta := viper.GetStringSlice("multi.sta")
	for idx, value := range newSta {
		fmt.Printf("sta[%d]: %v\n", idx, value)
	}

	fruit := viper.GetStringSlice("fruit")
	fmt.Printf("fruit: %v\n", fruit)

	ports := viper.GetIntSlice("ports")
	fmt.Printf("ports: %v\n", ports)

	// 读取不存在的字段，字符串为空，数值为0
	bad := viper.GetString("bad")
	bad1 := viper.GetInt("my.bad")
	fmt.Printf("bad: [%v] bad1: [%v]\n", bad, bad1)

	logdir := viper.GetString("loginfo.log.dir")
	fmt.Printf("logdir: %v\n", logdir)

	// 多级对象
	// tmpMap := make([]mapUrl_t, 0, 20)
	var tmpMap []mapUrl_t

	viper.UnmarshalKey("mymap.map_data", &tmpMap)

	for _, item := range tmpMap {
		fmt.Printf("name: %v url: %v\n", item.Name, item.Url)
	}

	// viper.WatchConfig()
	// viper.OnConfigChange(func(e fsnotify.Event) {
	// 	fmt.Println("配置发生变更：", e.Name)
	// })
}
