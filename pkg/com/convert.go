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

// 各种数字、字符等的格式转换

package com

import (
	"fmt"
    "bytes"
	"strconv"
	"math"
    "os"
    "io"
    "reflect"
	"encoding/hex"
    "encoding/gob"
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

func ToFixed1(value, scale int) (s string) {
    realValue := float64(value) / math.Pow10(scale)
    s = fmt.Sprintf("%v", realValue)
    return s
    //fmt.Printf("%v %v\r\n", realValue, str);
}

// TODO:有无更好的方式？
func ToFixed(value interface{}, scale int) (s string) {
	var realValue float64
	switch v := value.(type) {
	case float32:
        realValue = float64(v)
	case float64:
        realValue = float64(v)
	case int8:
        realValue = float64(v)
	case int16:
        realValue = float64(v)
	case int32:
        realValue = float64(v)
	case int64:
        realValue = float64(v)
	case uint8:
        realValue = float64(v)
	case uint16:
        realValue = float64(v)
	case uint32:
        realValue = float64(v)
	case uint64:
	    realValue = float64(v)
    case uint:
        realValue = float64(v)
	case int:
		realValue = float64(v)
	default:
		realValue = 0
	}
	realValue = realValue / math.Pow10(scale);
    s = fmt.Sprintf("%v", realValue);
	return s
}

func Round(value interface{}, scale int) (s string) {
    var realValue float64
	switch v := value.(type) {
	case float32:
        realValue = float64(v)
	case float64:
        realValue = float64(v)
	case int8:
        realValue = float64(v)
	case int16:
        realValue = float64(v)
	case int32:
        realValue = float64(v)
	case int64:
        realValue = float64(v)
	case uint8:
        realValue = float64(v)
	case uint16:
        realValue = float64(v)
	case uint32:
        realValue = float64(v)
	case uint64:
	    realValue = float64(v)
    case uint:
        realValue = float64(v)
	case int:
		realValue = float64(v)
	default:
		realValue = 0
	}
	realValue = realValue * math.Pow10(scale);
    s = fmt.Sprintf("%v", realValue);
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

// Marshal converts the given struct to a byte slice.
func Marshal(o interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	err := encoder.Encode(o)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Unmarshal reads from the byte slice `data` and decodes the struct into `o`.
func Unmarshal(data []byte, o interface{}) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)

	err := decoder.Decode(o)
	if err != nil {
		return err
	}

	return nil
}

func Dump(by []byte, len int) {
    line := 16
    n := len / line
    if len % line != 0 {
        n++
    }
    for i := 0; i < n; i++ {
        fmt.Printf("%08x  ", i*line)
        for j := 0; j < line; j++ {
            if i*line+j < len {
                fmt.Printf("%02x ", by[i*line+j])
            } else {
                fmt.Printf("   ")
            }
            if j == 7 {
                fmt.Printf(" ")
            }
        }
        fmt.Printf(" |")
        for j := 0; j<line && (i*line+j)<len; j++ {
            if (i*line+j) < len {
                c := by[i*line+j]
                if c >= ' ' && c < '~'{
                    fmt.Printf("%c", c)
                } else {
                    fmt.Printf(".")
                }
            } else {
                fmt.Printf("   ")
            }
        }
        fmt.Printf("|\n")
    }
}


func Dump1(w io.Writer, by []byte, len int) {
    line := 16
    n := len / line
    if len % line != 0 {
        n++
    }
    for i := 0; i < n; i++ {
        fmt.Fprintf(w, "%08x  ", i*line)
        for j := 0; j < line; j++ {
            if i*line+j < len {
                fmt.Fprintf(w, "%02x ", by[i*line+j])
            } else {
                fmt.Fprintf(w, "   ")
            }
            if j == 7 {
                fmt.Fprintf(w, " ")
            }
        }
        fmt.Fprintf(w, " |")
        for j := 0; j<line && (i*line+j)<len; j++ {
            if (i*line+j) < len {
                c := by[i*line+j]
                if c >= ' ' && c < '~'{
                    fmt.Fprintf(w, "%c", c)
                } else {
                    fmt.Fprintf(w, ".")
                }
            } else {
                fmt.Fprintf(w, "   ")
            }
        }
        fmt.Fprintf(w, "|\n")
    }
}

// 将数组、map等，按行打印，默认fmt.Println是一行
func PrintByLine(w io.Writer, msg interface{}) {
	if w == os.Stderr {
		fmt.Fprintf(os.Stderr, "error: ")
	}
	t := reflect.TypeOf(msg)
	switch t.Kind() {
	case reflect.Slice, reflect.Array:
		fmt.Fprintf(w, "[\n")
		v := reflect.ValueOf(msg)
		for i := 0; i < v.Len(); i++ {
			fmt.Fprintf(w, "  %v\n", v.Index(i))
		}
		fmt.Fprintf(w, "]\n")
	case reflect.Map:
		fmt.Fprintf(w, "[\n")
		v := reflect.ValueOf(msg)
		iter := v.MapRange()
		for iter.Next() {
			fmt.Fprintf(w, "  %v: %v\n", iter.Key(), iter.Value())
		}
		fmt.Fprintf(w, "]\n")
	default:
		fmt.Fprintf(w, "%v\n", msg)
	}
}