/*

<-chan signal 只接收，一般用于函数返回
out chan <- chan  只发送，一般作函数参数
*/
package test

import (
	"log"
	"sync"
	"testing"
	"time"
)

var g_download = make(chan bool)
var g_unzip = make(chan bool)
var g_copy = make(chan bool)
var g_doit = make(chan bool)
var g_cnt = 1

func downloadFromFTP() {
	println(" wait download...  ")
	<-g_download

	g_unzip <- true
}

func unzipFile() {
	for {
		println(" wait update... ")
		<-g_unzip
		println(" update ")
		//time.Sleep(2 * time.Second)
		g_copy <- true
		//g_doit <- true
	}
}

func copyFile() {
	for {
		println(" wait upload... ")
		<-g_copy
		//time.Sleep(3 * time.Second)
		println(" upload ")
		g_doit <- true
	}
}

func doit() {
	for {
		<-g_doit
		println(" doit!!!!!!!!!!!!!!!!... ", g_cnt)
		g_cnt += 1
		time.Sleep(3 * time.Second)
		g_doit <- true // 在这里再触发，似乎不行
	}

}

func TestChannelTask(t *testing.T) {
	cnt := 0
	println("start main")

	go unzipFile()

	go copyFile()

	go doit()

	go downloadFromFTP()

	for {
		// println(" . ")
		time.Sleep(1 * time.Second)
		if cnt == 3 {
			g_download <- true
		}
		cnt++
	}
	println("end main")
}

//////////////////////////////

type signal struct{}

func worker() {
	log.Println("worker is working")
	time.Sleep(1 * time.Second)
}

func spawn(f func()) <-chan signal {
	c := make(chan signal)
	go func() {
		log.Println("start worker")
		f()
		c <- signal(struct{}{})
	}()

	return c
}

func TestChannelTask2(t *testing.T) {
	log.Println("start worker in main")
	c := spawn(worker) //结束后返回 chan，主线程中等待之
	<-c
	log.Println("worker done in main")
}

/////////////////////////////
// 一个chan通知多个协程

func worker3(i int) {
	log.Printf("worker %v is working\n", i)
	time.Sleep(1 * time.Second)
	log.Printf("worker %v done\n", i)
}

func spawnGroup(f func(i int), num int, groupSignal <-chan signal) <-chan signal {
	c := make(chan signal)
	var wg sync.WaitGroup

	for i := 0; i < num; i++ {
		wg.Add(1)
		go func(i int) {
			<-groupSignal
			log.Printf("worker %v: start to working....\n", i)
			f(i)
			wg.Done()
		}(i)
	}

	// 用wg等待所有协程完成，再返回chan
	go func() {
		wg.Wait()
		c <- signal(struct{}{})
	}()

	return c
}

func TestChannelTask3(t *testing.T) {
	log.Println("test of group worker in main")

	groupSignal := make(chan signal)
	c := spawnGroup(worker3, 5, groupSignal)
	time.Sleep(3 * time.Second)
	log.Println("start group worker in main")
	close(groupSignal) // 发信号
	<-c
	log.Println("worker done in main")
}

////////////////////////////////////////////
func heartbeat() {

	heartbeat := time.NewTicker(3 * time.Second) // 3秒定时
	defer heartbeat.Stop()

	for {
		select {
		// case <-c:
		case <-heartbeat.C:
			log.Println(".")
		}
	}
}
func TestChannelTask4(t *testing.T) {
	log.Println("test hearbeat")
	heartbeat()
}
