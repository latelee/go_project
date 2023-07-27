/*
机场模拟并发  并行方案

理解：3条通道，每条通道，一个工作人员做所有的事


并发不是并行， 并行关乎结构，并行关乎执行
	-- Rob Pike，Go语言之父
*/
package test

import (
	"fmt"
	"testing"
	"time"
)

// const (
// 	checkIdTime   = 60
// 	checkBodyTime = 120
// 	checkXrayTime = 180
// )

func checkId1(id int) int {
	t := checkIdTime
	name := "check id"
	time.Sleep(time.Millisecond * time.Duration(t))
	fmt.Printf("\tgo-%v %v ok\n", id, name)
	return t
}

func checkBody1(id int) int {
	t := checkBodyTime
	name := "check body"
	time.Sleep(time.Millisecond * time.Duration(t))
	fmt.Printf("\tgo-%v %v ok\n", id, name)
	return t
}

func checkXray1(id int) int {
	t := checkXrayTime
	name := "check xray"
	time.Sleep(time.Millisecond * time.Duration(t))
	fmt.Printf("\tgo-%v %v ok\n", id, name)
	return t
}

func airportCheck1(id int) int {
	fmt.Printf("go-%v airportCheck ...\n", id)
	ret := 0

	ret += checkId1(id)
	ret += checkBody1(id)
	ret += checkXray1(id)

	fmt.Printf("go-%v airportCheck ok\n", id)

	return ret
}

func start(id int, f func(int) int, queue <-chan struct{}) <-chan int {
	c := make(chan int)
	go func() {
		total := 0
		for {
			_, ok := <-queue
			if !ok {
				c <- total
				return
			}
			total += f(id)
		}
	}()
	return c
}

func max(args ...int) int {
	n := 0
	for _, v := range args {
		if v > n {
			n = v
		}
	}
	return n
}

// 并行3通道执行，耗时：total time cost:  3600 (3协程)
func TestAirportCheck1(t *testing.T) {

	return
	total := 0

	c := make(chan struct{})
	c1 := start(1, airportCheck1, c)
	c2 := start(2, airportCheck1, c)
	c3 := start(3, airportCheck1, c)

	for i := 0; i < passengers; i++ {
		c <- struct{}{}
	}
	close(c)
	total = max(<-c1, <-c2, <-c3)
	fmt.Println("total time cost1: ", total)

}
