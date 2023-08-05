package test

import (
	"log"
	"runtime"
	"testing"
	"time"
)

func add(a, b int) int {
	return a + b
}

func deadloop() {
	for {
		add(1, 2)
	}
}

func TestDeadLock1(t *testing.T) {
	runtime.GOMAXPROCS(1)
	go deadloop()
	for {
		time.Sleep(time.Second * 1)
		log.Println("I got scheduled!")
	}
}

//////////////////////
