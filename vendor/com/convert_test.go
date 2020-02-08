// Copyright 2014 com authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package com

import (
	"fmt"
	"testing"
)

func TestConvert(t *testing.T) {

	fmt.Println("conv ", ToStr(100.56))

	fmt.Println("Tofixed: ", ToFixed(65, 2))
	fmt.Println("Round: ", Round(65, 2))

	hex := ToHexByte("dd")
	fmt.Printf("hex1 %v \n", hex)
	hex = ToHexByte("deadbeef")
	fmt.Printf("hex2 %v \n", hex)

	str := ToHexString(hex)
	fmt.Printf("hexstr %v \n", str)
	
}