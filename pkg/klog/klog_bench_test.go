package klog

import (
	"fmt"
	"testing"
)

//"k8s.io/klog"

func BenchmarkKlog(b *testing.B) {
	fmt.Println("klog bench...")
	b.ResetTimer()

	log_dir := "log"
	log_prefix := "mytest"
	log_mod := "450102DB"
	log_level := 3
	log_size := 10 * 1024 * 1024 // 10 MiB
	show_type := 2               // 1: std 2:file 3:std+file
	log_type := 1                // 日志机制，0：旧机制  1：新形式
	log_time := 5
	end_type := 0
	Init_normal(log_dir, log_prefix, log_mod, log_level, log_size, show_type, log_type, end_type, log_time)

	Printf("[system init] |%v \n", "initing...")
	Warnf("[system init] |%v \n", "init failed")

	for i := 0; i < b.N; i++ {
		for j := 0; j < 10000; j++ {
			Infof("this333333 is %v %v----------------------------\n", i, j)
			Warnf("this333333 is warn %v %v ----------------------------\n", i, j)
		}
	}
	FlushLog()
}
