/*
定时器
*/

package test

import (
	"log"
	"sync"
	"testing"
	"time"
)

const conscnt = 5

func consume1(c <-chan bool) bool {
	timer := time.NewTimer(conscnt * time.Second)
	defer timer.Stop()

	select {
	case b := <-c: // 传参控制行为
		if b == false {
			log.Printf("recv false, continue")
			return true
		}
		log.Printf("recv true, return")
		return false
	case <-timer.C:
		log.Printf("timer expired")
		return true
	}
}

func pro_cons_test1() {
	c := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(2) // 2个协程
	// 生产者
	go func() {
		// 隔1秒发一信号
		for i := 0; i < conscnt; i++ {
			time.Sleep(1 * time.Second)
			c <- false
		}
		time.Sleep(1 * time.Second)
		c <- true
		wg.Done()
	}()

	// 消费者
	go func() {
		for {
			if b := consume1(c); !b { // 等待 为false表示结束了，返回
				wg.Done()
				return
			}
		}
	}()

	wg.Wait()
}

func TestProConsume1(t *testing.T) {
	pro_cons_test1()
}
