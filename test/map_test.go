/*
map测试

json：解析后保存为map，可以直接用知道的值做索引，有ok做判断
函数指针：

map sync.Map使用：
存储不一定按顺序
*/
package test

import (
	"fmt"
	"sync"
	"testing"
)

/////////////////////////
// 简单的map测试
// int string 类型 复合字面量赋值
var ColorMap = map[int]string{
	0: "无颜色/保留/默认",
	1: "蓝色",
	2: "黄色",
	3: "黑色",
	4: "白色",
	9: "其他",
}

func showMap1(info string, m map[int]string) {
	fmt.Printf("%v, len: %v\n", info, len(m))
	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	for k, v := range m {
		fmt.Printf("[%v, %v]\n", k, v)
	}
	fmt.Printf(">>>>>>>>>>>>>>>>>>>>>>>>>>>\n\n")
}

func changeMap(m map[int]string) {

	m[0] = "修改后的颜色"
}

func TestMapFunc1(t *testing.T) {
	// 使用make创建，未指定容量
	myColor := make(map[int]string)

	// 遍历赋值
	for k, v := range ColorMap {
		myColor[k] = v
	}

	showMap1("using make map", myColor)

	k := 100

	myColor[100] = "另一颜色" // 直接赋值

	delete(myColor, 9) // 删除

	changeMap(myColor) // 尝试在函数中修改存在的key 预期：可修改

	showMap1("using make map1", myColor)

	k = k + 1 // 此值控制是否在map中
	if v, ok := myColor[k]; ok {
		fmt.Printf("key [%v] found: [%v]\n", k, v)
	} else {
		fmt.Printf("key [%v] NOT FOUND\n", k)
	}

}

/* 创建map性能对比

测试示例：
BenchmarkMapInit-8                   180           6270598 ns/op         5762590 B/op       4009 allocs/op
BenchmarkMapInitCap-8                309           3370926 ns/op         2833537 B/op       1680 allocs/op

*/
const mapSize = 100000

func BenchmarkMapInit(b *testing.B) {
	for n := 0; n < b.N; n++ {
		m := make(map[int]int)
		for i := 0; i < mapSize; i++ {
			m[i] = i
		}
	}
}

func BenchmarkMapInitCap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		m := make(map[int]int, mapSize)
		for i := 0; i < mapSize; i++ {
			m[i] = i
		}
	}
}

//////////////////////////////////////
// string ptr 类型 简单定义map的函数指针
var FuncMap = map[string]func(int, string){
	"111": map1,
	"222": map2,
}

func map1(a int, b string) {
	fmt.Println("111", a, b)
}

func map2(a int, b string) {
	fmt.Println("222", a, b)
}

func funcMapTest(id string) {
	a := 250
	b := "25.250"
	if fn, ok := FuncMap[id]; ok {
		fn(a, b)
	} else {
		fmt.Printf("id [%v] not found func ptr\n", id)
	}
}

func TestMapFunc2(t *testing.T) {
	funcMapTest("111")
	funcMapTest("222")
	funcMapTest("333") // 不存在的
}

/////////////////////////
// sync.Map使用
func showMap(smap sync.Map) {
	// 遍历打印，不做死循环
	smap.Range(func(k, v interface{}) bool {
		v1 := v.(string)
		fmt.Printf("key: %v -> %v\n", k, v1)
		return true
	})
}

func TestMapSync(t *testing.T) {
	var myMap sync.Map
	// 存
	myMap.Store(1, "111")
	myMap.Store(2, "222")
	myMap.Store(3, "333")
	showMap(myMap)

	// 删除
	myMap.Delete(1)
	showMap(myMap)

	myMap.Store(3, "333_111") // 重复写同一个key，会更新为新的
	showMap(myMap)

	//Load 方法，获得value
	if v, ok := myMap.Load(2); ok {
		fmt.Println("Load ->", v)
	}
	//LoadOrStore方法，获取或者保存
	//参数是一对key：value，如果该key存在且没有被标记删除则返回原先的value（不更新）和true；
	// 不存在则store，返回该value 和false
	if vv, ok := myMap.LoadOrStore(1, "c"); !ok { // 前面删除1了，这里重新保存
		fmt.Println("save new", vv)
	}
	if vv, ok := myMap.LoadOrStore(2, "c"); ok { // 2 一直存在
		fmt.Println("exist", vv, ok)
	}
	showMap(myMap)
}

//////////////////////////////////
