// Copyright 2013 com authors
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
	"strings"
	"testing"
)

func TestString(t *testing.T) {

	fmt.Println("isletter: ", IsLetter('1'))
	fmt.Println("reverse ", Reverse("abcd中文"))

	fmt.Println("to snake ", ToSnakeCase("HttpServer"))
}

func TestExpand(t *testing.T) {
	match := map[string]string{
		"domain":    "gowalker.org",
		"subdomain": "github.com",
	}
	s := "http://{domain}/{subdomain}/{0}/{1}"
	sR := "http://gowalker.org/github.com/unknwon/gowalker"
	if Expand(s, match, "unknwon", "gowalker") != sR {
		t.Errorf("Expand:\n Expect => %s\n Got => %s\n", sR, s)
	}
}

func checkString(aaa_str, bbb_str string) bool {
	sameCnt := 0
	// 用此法对比不准确
	if len(aaa_str) == len(bbb_str) {
		for i := 0; i < len(bbb_str); i++ {
			if strings.EqualFold(string(aaa_str[i]), string(bbb_str[i])) {
				sameCnt++
			}
		}
	}
	if sameCnt >= len(bbb_str)-3 {
		return true
	}
	return false
}

func checkRune(aaa_str, bbb_str string) bool {
	sameCnt := 0
	// 如有中文，用rune类型
	aa_str := String2Rune(aaa_str)
	bb_str := String2Rune(bbb_str)
	if len(aa_str) == len(bb_str) {
		for i := 0; i < len(aa_str); i++ {
			if aa_str[i] == bb_str[i] {
				sameCnt++
			}
		}
	}
	if sameCnt >= len(bb_str)-3 {
		return true
	}
	return false
}

func TestStringNum(t *testing.T) {

	var a []string = []string{"岑溪450481", "岑溪450481", "岑溪450481", "岑溪450481"}
	var b []string = []string{"岑溪450481", "芩溪450481", "芩溪458481", "梧州450487"}

	for i := 0; i < len(a); i++ {
		fmt.Printf("%v string result: %v %v\n", i, checkString(a[i], b[i]), checkRune(a[i], b[i]))
	}

}
