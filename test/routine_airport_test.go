/*
机场模拟并发 顺序
1. ID身份检查
2. 人身检查
3. X光检查

理解：一条通道，一个工作人员做所有的事

并发不是并行， 并行关乎结构，并行关乎执行
	-- Rob Pike，Go语言之父
*/
package test

import (
	"fmt"
	"testing"
	"time"
)

const (
	checkIdTime   = 60
	checkBodyTime = 120
	checkXrayTime = 180
)

var passengers int = 30

func checkId() int {
	t := checkIdTime
	name := "check id"
	time.Sleep(time.Millisecond * time.Duration(t))
	fmt.Printf("\t%v ok\n", name)
	return t
}

func checkBody() int {
	t := checkBodyTime
	name := "check body"
	time.Sleep(time.Millisecond * time.Duration(t))
	fmt.Printf("\t%v ok\n", name)
	return t
}

func checkXray() int {
	t := checkXrayTime
	name := "check xray"
	time.Sleep(time.Millisecond * time.Duration(t))
	fmt.Printf("\t%v ok\n", name)
	return t
}

func airportCheck() int {
	fmt.Println("airportCheck ...")
	ret := 0

	ret += checkId()
	ret += checkBody()
	ret += checkXray()
	fmt.Println("airportCheck ok")

	return ret
}

// 顺序执行，耗时：total time cost:  10800
func TestAirportCheck(t *testing.T) {
	return
	total := 0

	for i := 0; i < passengers; i++ {
		total += airportCheck()
	}
	fmt.Println("total time cost: ", total)
}
