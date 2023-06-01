/*
定时器
*/

package test

import (
	"log"
	"testing"
	"time"
)

func create_timer_by_afterfunc() {
	_ = time.AfterFunc(1*time.Second, func() {
		log.Println("timer created by AfterFunc() fired!")
	})
}

func create_timer_by_newtimer() {
	timer := time.NewTimer(2 * time.Second)
	select {
	case <-timer.C:
		log.Println("timer created by NewTimer() fired!")
	}
}

func create_timer_by_after() {
	select {
	case <-time.After(2 * time.Second):
		log.Println("timer created by After() fired!")
	}
}

func TestTimer1(t *testing.T) {
	create_timer_by_afterfunc()
	create_timer_by_newtimer()
	create_timer_by_after()
}
