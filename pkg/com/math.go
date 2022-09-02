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
	"math"
)

// PowInt is int type of math.Pow function.
func PowInt(x int, y int) int {
	if y <= 0 {
		return 1
	} else {
		if y%2 == 0 {
			sqrt := PowInt(x, y/2)
			return sqrt * sqrt
		} else {
			return PowInt(x, y-1) * x
		}
	}
}

// 四舍五入 -10.6 = -10
// 先取整数，如果小数大于0.5，则+1
func RoundClassic(r float64) (ret int) {
	r, f := math.Modf(r) // 分别取整数、小数，符号相同
	if f >= 0.5 {
		ret = int(r + 1)
	} else {
		ret = int(r)
	}
	return

}

// 四舍五入 -10.6 = -11
func RoundClassic1(r float64) int {

	return int(math.Floor(r + 0.5))

}

// 上浮或下降百分比，保留2位小数
// 方法：先扩大100倍（即2位），再取整，最后除以100
// 有的语言（如delphi），floor结果与golang不同，加上0.001（即扩大100倍后有小数点）
func UpNumber(item, percent float64) (newitem float64) {
	newitem = math.Floor(item*(1+percent)*100+0.001) / 100.0
	// newrate2 := math.Floor(item*(1+0.10)*100) / 100.0
	return
}
