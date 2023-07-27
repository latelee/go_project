/*
数值统计分析相关接口

go get gonum.org/v1/gonum
*/

package utils

import (
	"fmt"
	"math"
	// "gonum.org/v1/gonum/floats"
	// "gonum.org/v1/gonum/stat"
)

// 均值
func Mean(samples []float64) float64 {
	n := len(samples)
	mean := 0.0
	for _, v := range samples {
		mean += v
	}
	mean /= float64(n)

	return mean
}

// 标准差 总体标准差
func StdDev(samples []float64) float64 {
	n := len(samples)
	mean := 0.0
	for _, v := range samples {
		mean += v
	}
	mean /= float64(n)

	variance := 0.0
	for _, v := range samples {
		variance += (v - mean) * (v - mean)
	}

	fmt.Println("lll variance ", variance)
	variance /= float64(n) // n-1 为样本方差
	fmt.Println("222 variance ", variance)

	return math.Sqrt(variance)
}

func MeanStdDev(samples []float64, deltanum ...int) (mean, std float64) {
	n := len(samples)
	// mean := 0.0
	for _, v := range samples {
		mean += v
	}
	mean /= float64(n)

	variance := 0.0
	for _, v := range samples {
		variance += (v - mean) * (v - mean)
	}

	fmt.Println("lll variance ", variance)
	variance /= float64(n) // n-1 为样本方差
	fmt.Println("222 variance ", variance)
	mynum := 1
	if len(deltanum) == 1 {
		mynum = deltanum[0]
	}
	std = math.Sqrt(variance) * float64(mynum)

	return
}

func func_mean() {
	// 创建一组样本数据
	numbers := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}

	mean := Mean(numbers)
	dev := StdDev(numbers)
	fmt.Println("Mean:", mean, dev, mean-2*dev, mean+2*dev)

	mean, dev = MeanStdDev(numbers, 2)

	fmt.Println("Mean:", mean, dev)

}
