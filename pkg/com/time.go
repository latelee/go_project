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
	"strconv"
	"strings"
	"time"
)

type DateTimeS struct {
	Year   int
	Month  int
	Day    int
	Hour   int
	Minute int
	Second int

	Tick uint64
}

// Format unix time int64 to string
func Date(ti int64, format string) string {
	t := time.Unix(int64(ti), 0)
	return DateT(t, format)
}

// Format unix time string to string
func DateS(ts string, format string) string {
	i, _ := strconv.ParseInt(ts, 10, 64)
	return Date(i, format)
}

// Format time.Time struct to string
// MM - month - 01
// M - month - 1, single bit
// DD - day - 02
// D - day 2
// YYYY - year - 2006
// YY - year - 06
// HH - 24 hours - 03
// H - 24 hours - 3
// hh - 12 hours - 03
// h - 12 hours - 3
// mm - minute - 04
// m - minute - 4
// ss - second - 05
// s - second = 5
// SSS - ms
func DateT(t time.Time, format string) string {
	res := strings.Replace(format, "MM", t.Format("01"), -1)
	res = strings.Replace(res, "M", t.Format("1"), -1)
	res = strings.Replace(res, "DD", t.Format("02"), -1)
	res = strings.Replace(res, "D", t.Format("2"), -1)
	res = strings.Replace(res, "YYYY", t.Format("2006"), -1)
	res = strings.Replace(res, "YY", t.Format("06"), -1)
	res = strings.Replace(res, "HH", fmt.Sprintf("%02d", t.Hour()), -1)
	res = strings.Replace(res, "H", fmt.Sprintf("%d", t.Hour()), -1)
	res = strings.Replace(res, "hh", t.Format("03"), -1)
	res = strings.Replace(res, "h", t.Format("3"), -1)
	res = strings.Replace(res, "mm", t.Format("04"), -1)
	res = strings.Replace(res, "m", t.Format("4"), -1)
	res = strings.Replace(res, "SSS", fmt.Sprintf("%03d", t.Nanosecond()/1000000), -1)
	res = strings.Replace(res, "ss", t.Format("05"), -1)
	res = strings.Replace(res, "s", t.Format("5"), -1)

	return res
}

// DateFormat pattern rules.
var datePatterns = []string{
	// year
	"Y", "2006", // A full numeric representation of a year, 4 digits   Examples: 1999 or 2003
	"y", "06", //A two digit representation of a year   Examples: 99 or 03

	// month
	"m", "01", // Numeric representation of a month, with leading zeros 01 through 12
	"n", "1", // Numeric representation of a month, without leading zeros   1 through 12
	"M", "Jan", // A short textual representation of a month, three letters Jan through Dec
	"F", "January", // A full textual representation of a month, such as January or March   January through December

	// day
	"d", "02", // Day of the month, 2 digits with leading zeros 01 to 31
	"j", "2", // Day of the month without leading zeros 1 to 31

	// week
	"D", "Mon", // A textual representation of a day, three letters Mon through Sun
	"l", "Monday", // A full textual representation of the day of the week  Sunday through Saturday

	// time
	"g", "3", // 12-hour format of an hour without leading zeros    1 through 12
	"G", "15", // 24-hour format of an hour without leading zeros   0 through 23
	"h", "03", // 12-hour format of an hour with leading zeros  01 through 12
	"H", "15", // 24-hour format of an hour with leading zeros  00 through 23

	"a", "pm", // Lowercase Ante meridiem and Post meridiem am or pm
	"A", "PM", // Uppercase Ante meridiem and Post meridiem AM or PM

	"i", "04", // Minutes with leading zeros    00 to 59
	"s", "05", // Seconds, with leading zeros   00 through 59

	// time zone
	"T", "MST",
	"P", "-07:00",
	"O", "-0700",

	// RFC 2822
	"r", time.RFC1123Z,
}

// Parse Date use PHP time format.
func DateParse(dateString, format string) (time.Time, error) {
	replacer := strings.NewReplacer(datePatterns...)
	format = replacer.Replace(format)
	return time.ParseInLocation(format, dateString, time.Local)
}

// 字符串形式的时间转成时间戳
func DateStr2Stamp(dateString, format string) int64 {
	replacer := strings.NewReplacer(datePatterns...)
	format = replacer.Replace(format)
	time, _ := time.ParseInLocation(format, dateString, time.Local)
	return time.Unix()
}

func DateTime2Stamp(time time.Time) int64 {
	return time.Unix()
}

func ParseDuration(now time.Time, s, fmt string) string {
	t, _ := time.ParseDuration(s)
	newTime := now.Add(t)
	return DateT(newTime, fmt)
}

func GetNowDateTime(fmt string) string {
	return DateT(time.Now(), fmt)
}

func Sleep(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

// func Sleep(ms time.Duration) {
// 	time.Sleep(ms * time.Millisecond)
// }

// 返回系统时间类型
func ParseDate3(date []byte) (time.Time, int) {
	if len(date) != 19 {
		return time.Time{}, 0
	} else {
		year := (((int(date[0])-'0')*10+int(date[1])-'0')*10+int(date[2])-'0')*10 + int(date[3]) - '0'
		month := time.Month((int(date[5])-'0')*10 + int(date[6]) - '0')
		day := (int(date[8])-'0')*10 + int(date[9]) - '0'
		hour := (int(date[11])-'0')*10 + int(date[12]) - '0'
		minute := (int(date[14])-'0')*10 + int(date[15]) - '0'
		second := (int(date[17])-'0')*10 + int(date[18]) - '0'

		t := time.Date(year, month, day, hour, minute, second, 0, time.Local)
		// fmt.Println("time: ", t.Unix())

		return t, 0
	}
}

// 返回时间结构体
func ParseDate4(date string) (s DateTimeS, ret int) {
	ret = 0
	if len(date) != 19 {
		ret = -1
	} else {
		s.Year = (((int(date[0])-'0')*10+int(date[1])-'0')*10+int(date[2])-'0')*10 + int(date[3]) - '0'
		s.Month = (int(date[5])-'0')*10 + int(date[6]) - '0'
		s.Day = (int(date[8])-'0')*10 + int(date[9]) - '0'
		s.Hour = (int(date[11])-'0')*10 + int(date[12]) - '0'
		s.Minute = (int(date[14])-'0')*10 + int(date[15]) - '0'
		s.Second = (int(date[17])-'0')*10 + int(date[18]) - '0'

		// t := time.Date(s.Year, time.Month(s.Month), s.Day, s.Hour, s.Minute, s.Second, 0, time.Local)
		// s.Tick = uint64(t.Unix())

	}
	return
}
