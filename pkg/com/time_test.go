
package com

import (
	"fmt"
	"time"
	"testing"
)

func TestTime(t *testing.T) {
	fmt.Println("time test: ", DateT(time.Now(), "YYYY-MM-DD HH:00:00"))
	fmt.Println("time test: ", DateT(time.Now(), "YYYY-MM-DD HH:mm:ss"))
	fmt.Println("now data: ", GetNowTime("YYYY-MM-DD"), GetNowTime("YYYY-MM-DD 00:00:00"))
	Sleep(1000)
	fmt.Println("now data: ", GetNowTime("YYYY-MM-DD HH:mm:ss"))
}
