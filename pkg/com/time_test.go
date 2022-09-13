package com

import (
	"fmt"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	fmt.Println("time test: ", DateT(time.Now(), "YYYY-MM-DD HH:00:00"))
	fmt.Println("time test: ", DateT(time.Now(), "YYYY-MM-DD HH:mm:ss"))
	fmt.Println("time test: ", GetNowDateTime("YYYY-MM-DDTHH:mm:ss:SSS"))
}
