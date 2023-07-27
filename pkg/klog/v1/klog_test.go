package klog

import (
	"fmt"
	"testing"
	"time"
    "os"
	//"k8s.io/klog"
)

func TestKlog(t *testing.T) {
	fmt.Println("klog test...")

    log_dir := "log"
	log_prefix := "mytest"
	log_level := 3
	log_size := 10 * 1024 * 1024 // 10 MiB
	log_type := 2 // 1: std 2:file 3:std+file
	log_time := 5
    Init_normal(log_dir, log_prefix, log_level, log_size, log_type, log_time)

	i := 0
	go func() {
		for {
			Infof("this is %v ----------------------------\n", i)
			fmt.Printf("this is %v\n", i)
			time.Sleep(time.Duration(1000) * time.Millisecond)
			
			i++

			// TODO 退出前但未写到文件的日志，如何处理？
			if i>7 {
				fmt.Printf("exit is %v\n", i)
				
				Infoln("exit...")
				os.Exit(0)
			}
		}
	}()

	for {
		time.Sleep(time.Duration(1000) * time.Millisecond)
	}

}

func TestKlogFile(t *testing.T) {
	fmt.Println("klog test...")

}
