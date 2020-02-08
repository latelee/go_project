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
	"strconv"
	"math"
	"encoding/hex"
)

// Convert string to specify type.
type StrTo string

func (f StrTo) Exist() bool {
	return string(f) != string(0x1E)
}

func (f StrTo) Uint8() (uint8, error) {
	v, err := strconv.ParseUint(f.String(), 10, 8)
	return uint8(v), err
}

func (f StrTo) Int() (int, error) {
	v, err := strconv.ParseInt(f.String(), 10, 0)
	return int(v), err
}

func (f StrTo) Int64() (int64, error) {
	v, err := strconv.ParseInt(f.String(), 10, 64)
	return int64(v), err
}

func (f StrTo) Float64() (float64, error) {
	v, err := strconv.ParseFloat(f.String(), 64)
	return float64(v), err
}

func (f StrTo) MustUint8() uint8 {
	v, _ := f.Uint8()
	return v
}

func (f StrTo) MustInt() int {
	v, _ := f.Int()
	return v
}

func (f StrTo) MustInt64() int64 {
	v, _ := f.Int64()
	return v
}

func (f StrTo) MustFloat64() float64 {
	v, _ := f.Float64()
	return v
}

func (f StrTo) String() string {
	if f.Exist() {
		return string(f)
	}
	return ""
}

// Convert any type to string.
func ToStr(value interface{}, args ...int) (s string) {
	switch v := value.(type) {
	case bool:
		s = strconv.FormatBool(v)
	case float32:
		s = strconv.FormatFloat(float64(v), 'f', argInt(args).Get(0, -1), argInt(args).Get(1, 32))
	case float64:
		s = strconv.FormatFloat(v, 'f', argInt(args).Get(0, -1), argInt(args).Get(1, 64))
	case int:
		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
	case int8:
		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
	case int16:
		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
	case int32:
		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
	case int64:
		s = strconv.FormatInt(v, argInt(args).Get(0, 10))
	case uint:
		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
	case uint8:
		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
	case uint16:
		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
	case uint32:
		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
	case uint64:
		s = strconv.FormatUint(v, argInt(args).Get(0, 10))
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		s = fmt.Sprintf("%v", v)
	}
	return s
}

type argInt []int

func (a argInt) Get(i int, args ...int) (r int) {
	if i >= 0 && i < len(a) {
		r = a[i]
	} else if len(args) > 0 {
		r = args[0]
	}
	return
}

// HexStr2int converts hex format string to decimal number.
func HexStr2int(hexStr string) (int, error) {
	num := 0
	length := len(hexStr)
	for i := 0; i < length; i++ {
		char := hexStr[length-i-1]
		factor := -1

		switch {
		case char >= '0' && char <= '9':
			factor = int(char) - '0'
		case char >= 'a' && char <= 'f':
			factor = int(char) - 'a' + 10
		default:
			return -1, fmt.Errorf("invalid hex: %s", string(char))
		}

		num += factor * PowInt(16, i)
	}
	return num, nil
}

// Int2HexStr converts decimal number to hex format string.
func Int2HexStr(num int) (hex string) {
	if num == 0 {
		return "0"
	}

	for num > 0 {
		r := num % 16

		c := "?"
		if r >= 0 && r <= 9 {
			c = string(r + '0')
		} else {
			c = string(r + 'a' - 10)
		}
		hex = c + hex
		num = num / 16
	}
	return hex
}

func ToFixed(value, scale int) (s string) {
    realValue := float64(value) / math.Pow10(scale)
    s = fmt.Sprintf("%v", realValue)
    return s
    //fmt.Printf("%v %v\r\n", realValue, str);
}

/*
// TODO:改为任意值
func ToFixed(value interface{}, scale int) (s string) {
	
	switch v := value.(type) {
	case float32:
	case float64:
	case int8:
	case int16:
	case int32:
	case int64:
	case uint8:
	case uint16:
	case uint32:
	case uint64:
	case uint:
	case int:
		fmt.Println("type int")
		realValue := float64(v) / math.Pow10(scale);
		s = fmt.Sprintf("%v", realValue);
	default:
		fmt.Println("type uuu")		
	}
	
	fmt.Println("ddddddddddd")
	return s
    
    //fmt.Printf("%v %v\r\n", realValue, str);
}
*/

func Round(value, scale int) (s string) {
    realValue := float64(value) * math.Pow10(scale)
	s = fmt.Sprintf("%v", realValue)
	return s
}

// 十六进制字符串转十六进制，输出为byte类型
func ToHexByte(str string) (ob []byte) {
    ob, _ = hex.DecodeString(str);
    
    return;
}

// 十六进制数组转十六进制字符串，输出为对应的字符
// 如 4c 77数组，将转换成4c77字符串，可保存到文件
func ToHexString(b []byte) (ostr string) {
    ostr = hex.EncodeToString(b);
    return;
}


// 简单解析版本
func ReadCP56Time2a(buf []byte) string {
    u1 := uint16(buf[0]);
    u2 := uint16(buf[1]);
    totalMillis := u2 << 8 | u1;

    var seconds = totalMillis / 1000;
    //millis = totalMillis % 1000; // 毫秒，不使用
    var minute = uint8(buf[2]) & 0x3f; // 低6比特
    var hour = uint8(buf[3]) & 0x1f;   // 低5比特为小时 有标准说低7位，23小时为10111，5比特即可
    var day = uint8(buf[4]) & 0x1f; // 低5比特为日期，高3位表示一周的第几天，但不使用到
    var month = uint8(buf[5]) & 0x0f;  // 低4比特有效
    var year = (uint16(buf[6]) & 0x7f) + 2000;       // 低7比特有效，存储的值以2000为基础，所以要加为2000

    // 年月日时分秒
    return fmt.Sprintf("%.4v%.2v%.2v%.2v%.2v%.2v", year, month, day, 
                                       hour, minute, seconds);
}

// TODO：设置CP56时间
/*
func WriteCP56Time2a(date) [] byte {
    // 获取当前时间—
    var fullyear = 2019;
    var month = 2;
    var dayOfWeek = 1;
    var dayOfMonth = 1;
    var hour = 1;
    var minute = 1;
    var millis = 1;
    var totalMillis = seconds * 1000 + millis;
    var day = (dayOfWeek << 5) | dayOfMonth;

    var buf = [7]byte;
    buf[0] = 1;
    buf[1] = 1;
    buf[2] = minute;
    buf[3] = hour;
    buf[4] = day;
    buf[5] = month;
    buf[6] = fullyear - 2000;
}
*/