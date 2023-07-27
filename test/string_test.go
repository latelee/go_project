package test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

/*

字符串组装性能测试
结论：传值和传指针，性能相差不大
go test -bench . -benchmem xxx.go
*/
var s1 string = `GOdafdsfdsafg ldafiasdfh3a adlgiajdsf39fjsdvhcbzxcvadjfsadhgasdigsdghl i3dsfsadfhsfasdlfiasdfhsaglsadigsdfsadfasdfasdfadsf3`

func handleString(s string) {
	_ = s + "hello"
}

func handlePtrToString(s *string) {
	_ = *s + "hello"
}

func BenchmarkHandleString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		handleString(s1)
	}
}

func BenchmarkHandlePtrToString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		handlePtrToString(&s1)
	}
}

///////////////////////////////////////

/*
字符串构造
测试：
go test -bench Benchmark_string_*  -benchmem .\string_test.go

*/

var s2 []string = []string{
	"aaaaaaaaaaa ",
	"bbbbbbbbbbb ",
	"ccccccccccc ",
}

func makeString_plus(s []string) string {
	var ret string
	for _, v := range s {
		ret += v
	}
	return ret
}

// 多种变量的格式化，适用此法
func makeString_sprintf(s []string) string {
	var ret string
	for _, v := range s {
		ret = fmt.Sprintf("%s%s", ret, v)
	}
	return ret
}

// 如果是数组形式，适用此法
func makeString_join(s []string) string {
	return strings.Join(s, "")
}

func makeString_build(s []string) string {
	var b strings.Builder
	for _, v := range s {
		b.WriteString(v)
	}
	return b.String()
}

// 有初始化大小，下同
func makeString_buildinit(s []string) string {
	var b strings.Builder
	b.Grow(128)
	for _, v := range s {
		b.WriteString(v)
	}
	return b.String()
}

func makeString_buffer(s []string) string {
	var b bytes.Buffer
	for _, v := range s {
		b.WriteString(v)
	}
	return b.String()
}

func makeString_bufferinit(s []string) string {
	buf := make([]byte, 0, 128)
	b := bytes.NewBuffer(buf)
	for _, v := range s {
		b.WriteString(v)
	}
	return b.String()
}

func Benchmark_string_plus(b *testing.B) {
	for n := 0; n < b.N; n++ {
		makeString_plus(s2)
	}
}

func Benchmark_string_sprintf(b *testing.B) {
	for n := 0; n < b.N; n++ {
		makeString_sprintf(s2)
	}
}

func Benchmark_string_join(b *testing.B) {
	for n := 0; n < b.N; n++ {
		makeString_join(s2)
	}
}

func Benchmark_string_build(b *testing.B) {
	for n := 0; n < b.N; n++ {
		makeString_build(s2)
	}
}

func Benchmark_string_buildinit(b *testing.B) {
	for n := 0; n < b.N; n++ {
		makeString_buildinit(s2)
	}
}

func Benchmark_string_buffer(b *testing.B) {
	for n := 0; n < b.N; n++ {
		makeString_buffer(s2)
	}
}

func Benchmark_string_bufferinit(b *testing.B) {
	for n := 0; n < b.N; n++ {
		makeString_bufferinit(s2)
	}
}
