/*
定时器

同例2，只用一个定时器 定时器由外部传入，但协程延时时间过长
解决问题
*/

package test

import (
	"log"
	"sync"
	"testing"
	"time"
)

func consume3(c <-chan bool, timer *time.Timer) bool {
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}

	timer.Reset(conscnt * time.Second)

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

func pro_cons_test3() {
	timecnt := 7
	c := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(2) // 2个协程
	// 生产者
	go func() {
		// 隔1秒发一信号
		for i := 0; i < conscnt; i++ {
			time.Sleep(time.Duration(timecnt) * time.Second)
			c <- false
		}
		time.Sleep(time.Duration(timecnt) * time.Second)
		c <- true
		wg.Done()
	}()

	// 消费者
	go func() {
		timer := time.NewTimer(conscnt * time.Second) // 定时器在此定义
		for {
			if b := consume3(c, timer); !b { // 等待 为false表示结束了，返回
				wg.Done()
				return
			}
		}
	}()

	wg.Wait()
}

func TestProConsume3(t *testing.T) {
	pro_cons_test3()
}
