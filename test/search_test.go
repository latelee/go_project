/*
定时器

同例2，只用一个定时器 定时器由外部传入，但协程延时时间过长
解决问题
*/

package test

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type MyBaseInfo_t struct {
	Version string
	Id      string
	Name    string
	Name1   string
}

type ExtraInfo_t struct {
	Id    string
	Name  string
	Name1 string
}

// 针对大量数据数组，给出指定数量及索引，一般用于多协程并发处理
func MakeRoutineRange(total, num int) (start, end []int) {

	slice := total / num
	start = make([]int, num, num)
	end = make([]int, num, num)
	idxEnd := total
	for i := 0; i < num; i++ {
		start[i] = slice * i
		if i != num-1 {
			end[i] = (i + 1) * slice
		} else {
			end[i] = idxEnd
		}
	}

	return
}

func findItemFromId(arrays []MyBaseInfo_t, id string, atype int) (ret MyBaseInfo_t, ok bool) {
	atype = 1
	if atype == 0 {
		for i := 0; i < len(arrays); i++ {
			if arrays[i].Id == id {
				ret = arrays[i]
				ok = true
				return
			}
		}
	} else if atype == 1 {
		num := 10
		start, end := MakeRoutineRange(len(arrays), num)
		var wg sync.WaitGroup
		ch := make(chan struct{})
		wg.Add(num)
		// fmt.Println("len: ", len(start))
		for i := 0; i < len(start); i++ {
			go func(i, istart, iend int) {
				// fmt.Println("iii ", i, istart, iend)
				defer wg.Done()
				for i := istart; i < iend; i++ {
					if arrays[i].Id == id {
						ret = arrays[i]
						ok = true
						ch <- struct{}{}
						return
					}
				}
			}(i, start[i], end[i])
		}

		go func() {
			wg.Wait()
			ch <- struct{}{}
		}()

		<-ch

		// wg.Wait()
	}
	return
}

func paddingInfo(arrs []MyBaseInfo_t, extraInfo []ExtraInfo_t, atype int) (rets []MyBaseInfo_t) {
	rets = arrs
	if atype == 0 {
		for j := 0; j < len(extraInfo); j++ {
			item := extraInfo[j]
			if tmp, ok := findItemFromId(arrs, item.Id, atype); ok {
				rets[j].Name = tmp.Name
				rets[j].Name1 = tmp.Name1
			}
		}

		return
	}
	if atype == 1 {
		// 分10个协程 [start, end)  含start，不含end
		num := 10
		start, end := MakeRoutineRange(len(extraInfo), num)
		var wg sync.WaitGroup
		wg.Add(num)
		for i := 0; i < len(start); i++ {
			go func(istart, iend int) {
				defer wg.Done()
				// klog.Printf("%d: %v --> %v\n", i, istart, iend)
				for j := istart; j < iend; j++ {
					item := extraInfo[j]
					if tmp, ok := findItemFromId(arrs, item.Id, atype); ok {
						rets[j].Name = tmp.Name
						rets[j].Name1 = tmp.Name1
					}
				}
			}(start[i], end[i])
		}
		wg.Wait()
		return
	}

	return
}

////////////////////////
func makeData(len int) []MyBaseInfo_t {
	rets := make([]MyBaseInfo_t, 0, len)
	for i := 0; i < len; i++ {
		var tmp MyBaseInfo_t
		tmp.Id = fmt.Sprintf("id_%v", i)
		rets = append(rets, tmp)
	}
	return rets
}

func makePaddingData(len int) (rets []ExtraInfo_t) {
	rets = make([]ExtraInfo_t, 0, len)
	for i := 0; i < len; i++ {
		var tmp ExtraInfo_t
		tmp.Id = fmt.Sprintf("id_%v", i)
		tmp.Name = fmt.Sprintf("name_%v", 100+i)
		tmp.Name1 = fmt.Sprintf("name1_%v", 100+i)
		rets = append(rets, tmp)
	}
	return
}

func TestSearch(t *testing.T) {
	len := 10000
	baseInfo := makeData(len)
	padInfo := makePaddingData(len)
	t1 := time.Now()
	baseInfo = paddingInfo(baseInfo, padInfo, 1)

	fmt.Println("pad time: ", time.Since(t1))

}
