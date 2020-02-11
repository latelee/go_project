
package com

import (
	"fmt"
	"testing"
)

func TestConvert(t *testing.T) {

	fmt.Println("conv ", ToStr(100.56))

	fmt.Println("Tofixed: ", ToFixed(6.5, 2))
	fmt.Println("Round: ", Round(0.65, 4))

	hex := ToHexByte("dd")
	fmt.Printf("hex1 %v \n", hex)
	hex = ToHexByte("deadbeef")
	fmt.Printf("hex2 %v \n", hex)

	str := ToHexString(hex)
	fmt.Printf("hexstr %v \n", str)
	
}