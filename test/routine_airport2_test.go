/*
机场模拟并发  并发方案

理解：3条通道，每条通道，3个工作人员做3个事

并发不是并行， 并行关乎结构，并行关乎执行
	-- Rob Pike，Go语言之父
*/
package test

import (
	"fmt"
	"testing"
	"time"
)

func start2(id int, f func(int) int, next chan<- struct{}) (chan<- struct{},
	chan<- struct{}, <-chan int) {
	// 这里可能是指处理的队列资源（与乘客数量无关） 经测试，不是越大越好，也不是越小越好
	queue := make(chan struct{}, 10)
	quit := make(chan struct{})
	result := make(chan int)

	go func() {
		total := 0
		for {
			select {
			case <-quit:
				result <- total
				// fmt.Printf("\t+++ %v: time: %v\n", id, total)
				return
			case v := <-queue:
				total += f(id)
				// fmt.Printf("\t... %v: time: %v\n", id, total)
				if next != nil {
					next <- v
				}
			}
		}
	}()
	return queue, quit, result
}

func airportCheck2(id int, queue <-chan struct{}) {
	go func(id int) {
		fmt.Printf("go-channel%v channel is ready ...\n", id)
		// 按顺序，所以最后一步先做
		queue3, quit3, result3 := start2(id, checkXray1, nil)
		queue2, quit2, result2 := start2(id, checkBody1, queue3)
		queue1, quit1, result1 := start2(id, checkId1, queue2)

		for {
			select {
			case v, ok := <-queue:
				if !ok {
					close(quit1)
					close(quit2)
					close(quit3)
					total := max(<-result1, <-result2, <-result3)
					fmt.Printf("go-channel%v, check channel time: %v\n", id, total)
					fmt.Printf("go-channel%v, check channel closed\n", id)
					return
				}
				queue1 <- v
			}
		}
	}(id)
}

/*
并发，最大耗时2160
go-channel2, check channel time: 2160
go-channel2, check channel closed
go-channel1, check channel time: 2160
go-channel1, check channel closed
go-channel3, check channel time: 1080
go-channel3, check channel closed
*/
func TestAirportCheck2(t *testing.T) {
	queue := make(chan struct{}, passengers) // 按乘客数量开最大缓冲
	airportCheck2(1, queue)
	airportCheck2(2, queue)
	airportCheck2(3, queue)

	time.Sleep(1 * time.Second) // 确认上面的协程开启了
	for i := 0; i < passengers; i++ {
		queue <- struct{}{}
	}
	time.Sleep(5 * time.Second)
	close(queue)
	time.Sleep(10 * time.Second) // 延时足够的时间，防止提前退出
}
