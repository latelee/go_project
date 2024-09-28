package klog

import (
	"fmt"
	"os"
	"testing"
	"time"
)

//"k8s.io/klog"

func TestKlog(t *testing.T) {
	fmt.Println("klog test...")

	log_dir := "log"
	log_prefix := "450102DB_fee_fee-server1"
	log_mod := "450102DB"
	log_level := 3
	log_size := 10 * 1024 * 1024 // 10 MiB
	show_type := 2               // 1: std 2:file 3:std+file
	log_type := 1                // 日志机制，0：旧机制  1：新形式
	log_time := 5
	end_type := 0
	Init_normal(log_dir, log_prefix, log_mod, log_level, log_size, show_type, log_type, end_type, log_time)

	Infoln("init.....")

	Infoln("[系统初始化]", "init again.....")

	Printf("[system init] |%v \n", "initing...")
	Warnf("[system init] |%v \n", "init failed")

	i := 0
	go func() {
		for {
			Infof("this is %v ----------------------------\r\n", i)
			Warnf("this is warn %v ----------------------------\r\n", i)

			fmt.Printf("this is %v\n", i)
			time.Sleep(time.Duration(2000) * time.Millisecond)

			i++

			// TODO 退出前但未写到文件的日志，如何处理？
			if false && i > 7 {
				fmt.Printf("exit is %v\n", i)

				Infoln("exit...")
				Flush()
				os.Exit(0)

			}
		}
	}()

	for {
		time.Sleep(time.Duration(1000) * time.Millisecond)
	}

}

func TestKlog1(t *testing.T) {
	fmt.Println("klog bench...")

	log_dir := "log"
	log_prefix := "mytest"
	log_mod := "450102DB"
	log_level := 3
	log_size := 10 * 1024 * 1024 // 10 MiB
	show_type := 2               // 1: std 2:file 3:std+file
	log_type := 1                // 日志机制，0：旧机制  1：新形式
	log_time := 5
	end_type := 1
	Init_normal(log_dir, log_prefix, log_mod, log_level, log_size, show_type, log_type, end_type, log_time)

	Printf("[system init] |%v \n", "initing...")
	Warnf("[system init] |%v \n", "init failed")

	for i := 0; i < 10; i++ {
		t1 := time.Now()
		for j := 0; j < 100; j++ {
			Infof("this333333 is  %v----------------------------\n", j)
			Warnf("this333333 is warn %v ----------------------------\n", j)
		}

		fmt.Println("time ", time.Since(t1))
	}
	FlushLog()
}
