package klog

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// MaxSize is the maximum size of a log file in bytes.
var MaxSize uint64 = 1024 * 1024 * 10 // 1024 * 1024 * 1800

func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func isFile(filePath string) bool {
	f, e := os.Stat(filePath)
	if e != nil {
		return false
	}
	return !f.IsDir()
}

func mkDir(destPath string) error {
	if !isExist(destPath) {
		return os.MkdirAll(destPath, os.ModePerm)
	}
	return nil
}

func fileSize(file string) (uint64, error) {
	f, err := os.Stat(file)
	if err != nil {
		return 0, err
	}
	return uint64(f.Size()), nil
}

// 按文件名排序，可扩展至文件时间
type byName []os.FileInfo

//func (f byName) Less(i, j int) bool { return f[i].Name() < f[j].Name() } // 文件名升序，默认方式
func (f byName) Less(i, j int) bool { return f[i].Name() > f[j].Name() } // 文件名倒序
func (f byName) Len() int           { return len(f) }
func (f byName) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

// 按带数字的文件名排序
type byNumericalFilename []os.FileInfo

func (f byNumericalFilename) Len() int      { return len(f) }
func (f byNumericalFilename) Swap(i, j int) { f[i], f[j] = f[j], f[i] }

func (f byNumericalFilename) Less(i, j int) bool {
	pathA := f[i].Name()
	pathB := f[j].Name()

	// !! 根据需求，文件最后是数字，按其值降序排序，示例：ddd.log.x.1 ddd.log.x.2
	// 如有其它者，也可以修改
	a, err1 := strconv.Atoi(pathA[strings.LastIndex(pathA, ".")+1:])
	b, err2 := strconv.Atoi(pathB[strings.LastIndex(pathB, ".")+1:])

	// 整体文件（不含后缀名）名称排序,名称是数字
	// a, err1 := strconv.Atoi(pathA[0:strings.LastIndex(pathA, ".")])
	// b, err2 := strconv.Atoi(pathB[0:strings.LastIndex(pathB, ".")])

	// fmt.Println("---------- ", pathA, pathB, strings.LastIndex(pathA, "."), strings.LastIndex(pathB, "."), a, b)

	// 有错误，默认降序
	if err1 != nil || err2 != nil {
		return pathA > pathB
	}

	// 按数字降序
	return a > b
}

func getFileListByPrefix(dirPath, suffix string, needDir bool, num int) ([]string, error) {
	if !isExist(dirPath) {
		return nil, fmt.Errorf("given path does not exist: %s", dirPath)
	} else if isFile(dirPath) {
		return []string{dirPath}, nil
	}

	// Given path is a directory.
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}

	fis, err := dir.Readdir(0)
	if err != nil {
		return nil, err
	}

	fnum := len(fis)
	if fnum == 0 {
		return []string{}, nil
	}

	sort.Sort(byNumericalFilename(fis))

	if num == 0 {
		num = fnum
	} else {
		if num > fnum {
			num = fnum
		}
	}
	files := make([]string, 0, num)
	for i := 0; i < num; i++ {
		fi := fis[i]
		if strings.HasPrefix(fi.Name(), suffix) {
			if needDir {
				files = append(files, filepath.Join(dirPath, fi.Name()))
			} else {
				files = append(files, fi.Name())
			}
		}
	}

	return files, nil
}
