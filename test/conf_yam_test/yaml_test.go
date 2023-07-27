/*
go test -v -run TestYaml

输出：
test of yaml...
need: true name: conf file
version: 2
time: 2020-10-03T09:21:13
empty: nul
text: hello
world!

name: late \n lee, name1: late
 lee age: 99
sta[0]: 110 210 ddd 99
sta[1]: 133 135 1 2 1588 1509
sta[2]: 310-410
sta[3]: 333-444
fruit: [apple apple1 apple2 apple3 apple4 apple5]
bad: [] bad1: [0]
result true: [[1 1 1 1 1 1 1 1 1 1 1]]
result1 true: [[true true true true true true true true true true true]]
result false: [[0 0 0 0 0 0 0 0 0 0 0]]
result1 false: [[false false false false false false false false false false false]]
logdir: log
name: 在线 url: http://abc.com
name: 离线 url: http://ccc.com

*/

package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/viper"
)

var (
	cfgFile string
)

type mapUrl_t struct {
	Name  string `json:"name"`
	Attri string `json:"attri"`
	Url   string `json:"url"`
}

func TestYaml(t *testing.T) {
	fmt.Println("test of yaml...")

	// 设置配置文件的2种方式
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("./")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// 读取
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("'config.yaml' file read error:", err)
		os.Exit(0)
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

	// 读取不存在的字段，字符串为空，数值为0
	bad := viper.GetString("bad")
	bad1 := viper.GetInt("my.bad")
	fmt.Printf("bad: [%v] bad1: [%v]\n", bad, bad1)

	// 按数值、字符串读取on、off等值
	result := viper.GetIntSlice("result_true")
	fmt.Printf("result true: [%v]\n", result)
	result1 := viper.GetStringSlice("result_true")
	fmt.Printf("result1 true: [%v]\n", result1)

	result = viper.GetIntSlice("result_false")
	fmt.Printf("result false: [%v]\n", result)
	result1 = viper.GetStringSlice("result_false")
	fmt.Printf("result1 false: [%v]\n", result1)

	logdir := viper.GetString("loginfo.log.dir")
	fmt.Printf("logdir: %v\n", logdir)

	// 多级对象
	// tmpMap := make([]mapUrl_t, 0, 20)
	var tmpMap []mapUrl_t

	viper.UnmarshalKey("mymap.map_data", &tmpMap)

	for _, item := range tmpMap {
		fmt.Printf("name: %v url: %v\n", item.Name, item.Url)
	}
}
