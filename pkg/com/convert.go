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
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// Convert string to specify type.
type StrTo string

func (f StrTo) Exist() bool {
	return true
	// return string(f) != fmt.Sprint(0x1E) // 没明白，为何不等于0x1e（即30）
}

func (f StrTo) Uint8() uint8 {
	v, _ := strconv.ParseUint(f.String(), 10, 8)
	return uint8(v)
}

func (f StrTo) Int() int {
	v, _ := strconv.ParseInt(f.String(), 10, 0)
	return int(v)
}

func (f StrTo) Int64() int64 {
	v, _ := strconv.ParseInt(f.String(), 10, 64)
	return int64(v)
}

func (f StrTo) Float64() float64 {
	v, _ := strconv.ParseFloat(f.String(), 64)
	return float64(v)
}

func (f StrTo) Uint8Hex() uint8 {
	v, _ := strconv.ParseUint(f.String(), 16, 8)
	return uint8(v)
}

func (f StrTo) IntHex() int {
	v, _ := strconv.ParseInt(f.String(), 16, 0)
	return int(v)
}

func (f StrTo) Int64Hex() int64 {
	v, _ := strconv.ParseInt(f.String(), 16, 64)
	return int64(v)
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
func HexStr2int(hexStr string) int {
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
			return -1
		}

		num += factor * PowInt(16, i)
	}
	return num
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
			c = fmt.Sprint(r + '0')
		} else {
			c = fmt.Sprint(r + 'a' - 10)
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
	realValue = realValue / math.Pow10(scale)
	s = fmt.Sprintf("%v", realValue)
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
	realValue = realValue * math.Pow10(scale)
	s = fmt.Sprintf("%v", realValue)
	return s
}

// 十六进制字符串转十六进制，输出为byte类型
func ToHexByte(str string) (ob []byte) {
	ob, _ = hex.DecodeString(str)

	return
}

// 十六进制数组转十六进制字符串，输出为对应的字符
// 如 4c 77数组，将转换成4c77字符串，可保存到文件
func ToHexString(b []byte) (ostr string) {
	ostr = hex.EncodeToString(b)
	return
}

// TODO：移动到其它文件
// 简单解析版本
func ReadCP56Time2a(buf []byte) string {
	u1 := uint16(buf[0])
	u2 := uint16(buf[1])
	totalMillis := u2<<8 | u1

	var seconds = totalMillis / 1000
	//millis = totalMillis % 1000; // 毫秒，不使用
	var minute = uint8(buf[2]) & 0x3f         // 低6比特
	var hour = uint8(buf[3]) & 0x1f           // 低5比特为小时 有标准说低7位，23小时为10111，5比特即可
	var day = uint8(buf[4]) & 0x1f            // 低5比特为日期，高3位表示一周的第几天，但不使用到
	var month = uint8(buf[5]) & 0x0f          // 低4比特有效
	var year = (uint16(buf[6]) & 0x7f) + 2000 // 低7比特有效，存储的值以2000为基础，所以要加为2000

	// 年月日时分秒
	return fmt.Sprintf("%.4v%.2v%.2v%.2v%.2v%.2v", year, month, day,
		hour, minute, seconds)
}

// TODO：测试
func Str2Float64(str string, len int) float64 {
	lenstr := "%." + strconv.Itoa(len) + "f"
	value, _ := strconv.ParseFloat(str, 64)
	nstr := fmt.Sprintf(lenstr, value)
	val, _ := strconv.ParseFloat(nstr, 64)
	return val
}

func Str2Float64Round(str string, prec int, round bool) float64 {
	f, _ := strconv.ParseFloat(str, 64)
	return precision(f, prec, round)
}

func precision(f float64, prec int, round bool) float64 {
	pow10_n := math.Pow10(prec)
	if round {
		return math.Trunc((f+0.5/pow10_n)*pow10_n) / pow10_n
	}
	return math.Trunc((f)*pow10_n) / pow10_n
}

func Str2Int(str string) int {
	var s StrTo = StrTo(str)
	return s.Int()
}

func Str2Int64(str string) int64 {
	var s StrTo = StrTo(str)
	return s.Int64()
}

func Abs(value int) int {
	return int(math.Abs(float64(value)))
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

// convert map to struct
func Map2Struct(jmap interface{}, s interface{}) error {
	tmpDataJson, err := json.Marshal(jmap)
	if err != nil {
		return err
	}
	err = json.Unmarshal(tmpDataJson, &s)
	if err != nil {
		return err
	}
	return nil
}

func Dump(by []byte, len int) {
	line := 16
	n := len / line
	if len%line != 0 {
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
		for j := 0; j < line && (i*line+j) < len; j++ {
			if (i*line + j) < len {
				c := by[i*line+j]
				if c >= ' ' && c < '~' {
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
	if len%line != 0 {
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
		for j := 0; j < line && (i*line+j) < len; j++ {
			if (i*line + j) < len {
				c := by[i*line+j]
				if c >= ' ' && c < '~' {
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
// TODO map的测试
func PrintByLine(data interface{}) {
	var w io.Writer
	w = os.Stdout // set to stdout
	buffer, _ := PrintByLine2Buffer(data)
	fmt.Fprintf(w, "%v\n", buffer)
}

func PrintByLine2Buffer(data interface{}) (buffer string, num int) {
	buffer = ""

	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	ShowTitle := func(myref reflect.Value) (buffer string) {
		if myref.Type().Name() == "string" {
			return
		}
		for i := 0; i < myref.NumField(); i++ {
			buffer += fmt.Sprintf("%v ", myref.Type().Field(i).Name)
		}
		buffer += fmt.Sprintf("\n")
		return
	}

	switch t.Kind() {
	case reflect.Slice, reflect.Array:
		num = v.Len()
		title := ShowTitle(v.Index(0))
		buffer += fmt.Sprintf("%s", title)
		for i := 0; i < v.Len(); i++ {
			buffer += fmt.Sprintf("%d %v\n", i, v.Index(i))
		}
	case reflect.Map:
		num = v.Len()
		iter := v.MapRange()
		i := 0
		for iter.Next() {
			if i == 0 {
				title := ShowTitle(iter.Value())
				buffer += fmt.Sprintf("%s", title)
			}
			// 不加key
			// buffer += fmt.Sprintf("%d %v: %v\n", i+1, iter.Key(), iter.Value())
			buffer += fmt.Sprintf("%d %v\n", i, iter.Value())
			i += 1
		}
	default:
		// fmt.Fprintf(w, "%v\n", data)
		num = 1
		title := ShowTitle(v)
		buffer += fmt.Sprintf("%s", title)
		buffer += fmt.Sprintf("%v\n", v)
	}

	// fmt.Print(buffer)

	return
}

/*
参考
github.com/kataras/tablewriter
github.com/lensesio/tableprinter

*/

func indirectValue(val reflect.Value) reflect.Value {
	if kind := val.Kind(); kind == reflect.Interface || kind == reflect.Ptr {
		return val.Elem()
	}

	return val
}

func checkSkipNames(a string, b []string) bool {
	for _, item := range b {
		if item == a {
			return true
		}
	}
	return false
}

// 结构体的字段名称
func GetStructName(myref reflect.Value, title string, names []string, sep string) (buffer string) {
	// 注：有可能传递string数组，此时没有“标题”一说，返回
	if myref.Kind().String() != "struct" {
		return
	}
	// 如果不指定标题，则用结构体的
	if len(title) == 0 {
		for i := 0; i < myref.NumField(); i++ {
			if ok := checkSkipNames(myref.Type().Field(i).Name, names); ok {
				continue
			}
			buffer += fmt.Sprintf("%v %v ", sep, myref.Type().Field(i).Name)
		}
		buffer += fmt.Sprintf("%v\n", sep)
	} else {
		buffer += fmt.Sprintf("%s\n", title)
	}

	// 内部直接打印|---| 格式，不需要外部手动添加
	if sep == "|" {
		for i := 0; i < myref.NumField(); i++ {
			if ok := checkSkipNames(myref.Type().Field(i).Name, names); ok {
				continue
			}
			buffer += fmt.Sprintf("| --- ")
		}
		buffer += fmt.Sprintf("|\n")
	}
	return
}

// 将 | 替换为 <br>
// 将 , 替换为 |
func replaceI(text string, sep string) (ret string) {
	// 下面2种方法都可以
	// reg := regexp.MustCompile(`\|`)
	// ret = reg.ReplaceAllString(text, `${1}<br/>`)
	if sep == "|" {
		ret = strings.Replace(text, "|", "<br>", -1)
	} else if sep == "," {
		ret = strings.Replace(text, ",", "|", -1)
	}
	return ret
}

// 结构体的值
func GetStructValue(myref reflect.Value, names []string, sep string) (buffer string) {
	// fmt.Println("value kind", myref.Kind().String())
	// 注：有可能传递string数组，此时没有“字段”一说，返回原本的内容
	if myref.Kind().String() == "string" {
		return myref.Interface().(string)
	}
	if myref.Kind().String() != "struct" {
		return
	}

	for i := 0; i < myref.NumField(); i++ {
		if ok := checkSkipNames(myref.Type().Field(i).Name, names); ok {
			continue
		}
		// 判断是否包含|，有则替换，其必须是string类型，其它保持原有的
		t := myref.Field(i).Type().Name()
		if t == "string" {
			var str string = myref.Field(i).Interface().(string)
			str = replaceI(str, sep)
			buffer += fmt.Sprintf("%v %v ", sep, str)
		} else {
			buffer += fmt.Sprintf("%v %v ", sep, myref.Field(i).Interface())
		}
	}
	buffer += fmt.Sprintf("%v\n", sep)

	return
}

func PrintStructMarkdown(data interface{}, title string, skipNames ...string) {
	var w io.Writer
	w = os.Stdout // set to stdout
	buffer, num := ConvertStruct2Markdown(data, title, skipNames...)
	fmt.Fprintf(w, "total: %v  \n", num)
	fmt.Fprintf(w, "%v\n", buffer)
}

func ConvertStruct2Markdown(data interface{}, title string, skipNames ...string) (buffer string, num int) {
	return convertStruct(data, title, "|", skipNames...)
}

func ConvertStruct2CSV(data interface{}, title string, skipNames ...string) (buffer string, num int) {
	return convertStruct(data, title, ",", skipNames...)
}

/*
功能：指定结构体data，其可为slice map 单独结构体
	 指定自定义标题，为空则使用结构体字段
	 指定忽略的字段名称（即结构体字段的变量）
	 按结构体定义的顺序列出，如自定义标题，则必须保证一致。
*/
func convertStruct(data interface{}, title, sep string, skipNames ...string) (buffer string, num int) {
	buffer = ""

	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	var skipNamess []string
	for _, item := range skipNames {
		skipNamess = append(skipNamess, item)
	}

	// 打印结构体字段标志
	// innertitle := false
	printHead := false
	// if len(title) == 0 {
	// 	innertitle = true
	// }

	// 不同类型的，其索引方式不同，故一一判断使用
	switch t.Kind() {
	case reflect.Slice, reflect.Array:
		num = v.Len()
		buffer += GetStructName(v.Index(0), title, skipNamess, sep)
		for i := 0; i < v.Len(); i++ {
			buffer += GetStructValue(v.Index(i), skipNamess, sep)
		}
	case reflect.Map:
		num = v.Len()
		iter := v.MapRange()
		for iter.Next() {
			if !printHead {
				buffer += GetStructName(iter.Value(), title, skipNamess, sep)
				printHead = true
			}
			buffer += GetStructValue(iter.Value(), skipNamess, sep)
		}
	default:
		num = 1 // 单独结构体不能用Len，单独赋值
		if !printHead {
			buffer += GetStructName(v, title, skipNamess, sep)
			printHead = true
		}
		buffer += GetStructValue(v, skipNamess, sep)
	}

	return
}

/*
 将csv格式转成markdown
 整体看是string，但以\n分隔，再以,分隔
 如有标题，则为完整的markdown，否则只以|分隔
 buf, num = ConvertStruct2Markdown(objects, "| Name | Value | Size | Guard |\n| --- | --- | --- | ++++ |")
*/
func ConvertCSV2Markdown(data string, title string, skipNames ...string) (buffer string, num int) {
	buffer = ""
	if len(title) != 0 { // 标题
		// items := strings.Split(string(title), "\n")
		buffer += fmt.Sprintf("%s\n", title)
		items := strings.Split(title, "|")
		// fmt.Println("title len: ", len(items))
		for i := 0; i < len(items); i++ {
			if len(items[i]) == 0 { // 如果字段为0，不纳入计算
				continue
			}
			buffer += fmt.Sprintf("| --- ")
		}
		buffer += fmt.Sprintf("|\n")
	}

	items := strings.Split(data, "\n")
	for i := 1; i < len(items); i++ { // 第一行为默认的标题，不用
		item := items[i]
		if len(item) == 0 {
			continue
		}
		itemm := strings.Split(item, ",")
		// fmt.Println("len: ", itemm, len(itemm))
		for _, itemmm := range itemm {
			if len(itemmm) == 0 {
				continue
			}
			itemmm = replaceI(itemmm, "|")
			buffer += fmt.Sprintf("| %v ", itemmm)
		}
		buffer += fmt.Sprintf("|\n")
		if i > 10 {
			num = i
			break
		}
	}
	return
}

// 将结构体打印到文件，第一行为结构体变量 - 已不用
// func SaveByLine(filename string, structname interface{}, data interface{}) {
// 	w, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
// 	if err != nil {
// 		return
// 	}
// 	defer w.Close()

// 	t := reflect.TypeOf(data)
// 	v := reflect.ValueOf(data)
// 	if v.Len() == 0 {
// 		return
// 	}
// 	fmt.Fprintf(w, "total: %d\n", v.Len())

// 	var nop []string
// 	title := GetStructName(v, "", nop)
// 	fmt.Fprintf(w, "%v\n", title)

// 	//fmt.Fprintf(w, "[\n")
// 	switch t.Kind() {
// 	case reflect.Slice, reflect.Array:
// 		for i := 0; i < v.Len(); i++ {
// 			fmt.Fprintf(w, "%d  %v\n", i+1, v.Index(i))
// 		}
// 	case reflect.Map:
// 		iter := v.MapRange()
// 		i := 0
// 		for iter.Next() {
// 			fmt.Fprintf(w, "%d %v: %v\n", i+1, iter.Key(), iter.Value())
// 			i += 1
// 		}
// 	default:
// 		fmt.Fprintf(w, "%v\n", data)
// 	}

// }

//获取结构体中字段的名称及对应的类型
func GetStructFieldName(structName interface{}) []string {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		fmt.Println("Check type error not Struct")
		return nil
	}
	object := reflect.ValueOf(structName)
	myref := object.Elem()
	typeOfType := myref.Type()

	fieldNum := myref.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		result = append(result, typeOfType.Field(i).Name)
	}
	return result
}

func GetStructFieldType(structName interface{}) []string {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		fmt.Println("Check type error not Struct")
		return nil
	}
	object := reflect.ValueOf(structName)
	myref := object.Elem()

	fieldNum := myref.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		result = append(result, myref.Field(i).Type().Name())
	}
	return result
}

// 对比小数点后N位，N由point指定
func IsFloatEqual(f1, f2 float64, point ...int) bool {
	var mypoint int = 4
	if point != nil {
		mypoint = point[0]
	}
	if mypoint == 1 {
		return math.Abs(f1-f2) <= 0.1
	} else if mypoint == 2 {
		return math.Abs(f1-f2) <= 0.01
	} else if mypoint == 3 {
		return math.Abs(f1-f2) <= 0.001
	} else if mypoint == 4 {
		return math.Abs(f1-f2) <= 0.0001
	} else if mypoint == 5 {
		return math.Abs(f1-f2) <= 0.00001
	} else {
		return math.Abs(f1-f2) <= 0.0001
	}
}

/* 结构体对比，参数是指针类型
point指定浮点对比精度
names 为不对比的字段名称

遍历结构体，判断类型，其它可添加-->每个类型都要一一添加，是否有方便方式？
如不同，返回不同的字段信息
*/
func CompareStruct(data, data1 interface{}, point int, skipname ...string) (info string) {
	ret := false
	info = ""
	// 预防为空
	if IsInterfaceNil(data) || IsInterfaceNil(data1) {
		return
	}

	var skipNames []string
	for _, item := range skipname {
		skipNames = append(skipNames, item)
	}

	check := func(a string, b []string) bool {
		for _, item := range b {
			if strings.EqualFold(item, a) {
				return true
			}
		}
		return false
	}

	object := reflect.ValueOf(data)
	myref := object.Elem()
	//typeOfType := myref.Type()
	for i := 0; i < myref.NumField(); i++ {
		// 如果不需要对比，跳过
		if ok := check(myref.Type().Field(i).Name, skipNames); ok {
			continue
		}

		// 浮点数精度
		mypoint := 4
		if point != 0 {
			mypoint = point
		}
		field := myref.Field(i)
		field1 := reflect.ValueOf(data1).Elem().Field(i)
		switch field.Type().Name() {
		case "int8":
			ret = field.Interface().(int8) == field1.Interface().(int8)
		case "uint8":
			ret = field.Interface().(uint8) == field1.Interface().(uint8)
		case "int16":
			ret = field.Interface().(int16) == field1.Interface().(int16)
		case "uint16":
			ret = field.Interface().(uint16) == field1.Interface().(uint16)
		case "int":
			ret = field.Interface().(int) == field1.Interface().(int)
		case "uint":
			ret = field.Interface().(uint) == field1.Interface().(uint)
		case "int32":
			ret = field.Interface().(int32) == field1.Interface().(int32)
		case "uint32":
			ret = field.Interface().(uint32) == field1.Interface().(uint32)
		case "int64":
			ret = field.Interface().(int64) == field1.Interface().(int64)
		case "uint64":
			ret = field.Interface().(uint64) == field1.Interface().(uint64)
		case "string":
			ret = field.Interface().(string) == field1.Interface().(string)
		case "float64":
			ret = IsFloatEqual(field.Interface().(float64), field1.Interface().(float64), mypoint)
		case "float32":
			ret = IsFloatEqual(float64(field.Interface().(float32)), float64(field1.Interface().(float32)), mypoint)
		case "bool":
			ret = field.Interface().(bool) == field1.Interface().(bool)
		default:
			ret = false
		}
		if ret == false {
			info1 := fmt.Sprintf("%s: %v != %v ", myref.Type().Field(i).Name, field.Interface(), field1.Interface())
			info = info + info1
		}
	}
	return info
}

// 检查指定的结构体字段是否为空，仅限string类型
func CheckStructEmpty(data interface{}, checkname ...string) (ret int) {
	ret = 0

	// 预防为空
	if IsInterfaceNil(data) {
		ret = 1
		return
	}

	var checkNames []string
	for _, item := range checkname {
		checkNames = append(checkNames, item)
	}

	check := func(a string, b []string) bool {
		for _, item := range b {
			if strings.EqualFold(item, a) {
				return true
			}
		}
		return false
	}

	object := reflect.ValueOf(data)
	myref := object.Elem()
	// fmt.Printf("myref: %v %v\n", data, object)
	//typeOfType := myref.Type()
	for i := 0; i < myref.NumField(); i++ {
		// 如果不需要对比，跳过
		if ok := check(myref.Type().Field(i).Name, checkNames); !ok {
			continue
		}

		field := myref.Field(i)
		switch field.Type().Name() {
		case "string":
			ret1 := field.Interface().(string) == ""
			if ret1 {
				ret = 1
			}
		default:
			ret = 0
		}
	}
	return
}

// json格式化美观些
func FormatJson(jsonstr string) (ret string, err error) {
	var out bytes.Buffer
	err = json.Indent(&out, []byte(jsonstr), "", "    ")
	ret = out.String()
	return
}

func FormatJsonString(jsonstr string) (ret string) {
	var out bytes.Buffer
	_ = json.Indent(&out, []byte(jsonstr), "", "    ")
	ret = out.String()
	return
}

func PackJson(jsonstr string) string {
	jsonstr = strings.Replace(jsonstr, " ", "", -1)
	jsonstr = strings.Replace(jsonstr, "\t", "", -1)
	jsonstr = strings.Replace(jsonstr, "\r\n", "", -1)
	jsonstr = strings.Replace(jsonstr, "\n", "", -1)

	return jsonstr
}

func IsInterfaceNil(i interface{}) bool {
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}
