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
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// IsDir returns true if given path is a directory,
// or returns false when it's a file or does not exist.
func IsDir(dir string) bool {
	f, e := os.Stat(dir)
	if e != nil {
		return false
	}
	return f.IsDir()
}

func statDir(dirPath, recPath string, includeDir, isDirOnly, followSymlinks bool) ([]string, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	fis, err := dir.Readdir(0)
	if err != nil {
		return nil, err
	}

	statList := make([]string, 0)
	for _, fi := range fis {
		if strings.Contains(fi.Name(), ".DS_Store") {
			continue
		}

		relPath := path.Join(recPath, fi.Name())
		curPath := path.Join(dirPath, fi.Name())
		if fi.IsDir() {
			if includeDir {
				statList = append(statList, relPath+"/")
			}
			s, err := statDir(curPath, relPath, includeDir, isDirOnly, followSymlinks)
			if err != nil {
				return nil, err
			}
			statList = append(statList, s...)
		} else if !isDirOnly {
			statList = append(statList, relPath)
		} else if followSymlinks && fi.Mode()&os.ModeSymlink != 0 {
			link, err := os.Readlink(curPath)
			if err != nil {
				return nil, err
			}

			if IsDir(link) {
				if includeDir {
					statList = append(statList, relPath+"/")
				}
				s, err := statDir(curPath, relPath, includeDir, isDirOnly, followSymlinks)
				if err != nil {
					return nil, err
				}
				statList = append(statList, s...)
			}
		}
	}
	return statList, nil
}

// StatDir gathers information of given directory by depth-first.
// It returns slice of file list and includes subdirectories if enabled;
// it returns error and nil slice when error occurs in underlying functions,
// or given path is not a directory or does not exist.
//
// Slice does not include given path itself.
// If subdirectories is enabled, they will have suffix '/'.
func StatDir(rootPath string, includeDir ...bool) ([]string, error) {
	if !IsDir(rootPath) {
		return nil, errors.New("not a directory or does not exist: " + rootPath)
	}

	isIncludeDir := false
	if len(includeDir) >= 1 {
		isIncludeDir = includeDir[0]
	}
	return statDir(rootPath, "", isIncludeDir, false, false)
}

// LstatDir gathers information of given directory by depth-first.
// It returns slice of file list, follows symbolic links and includes subdirectories if enabled;
// it returns error and nil slice when error occurs in underlying functions,
// or given path is not a directory or does not exist.
//
// Slice does not include given path itself.
// If subdirectories is enabled, they will have suffix '/'.
func LstatDir(rootPath string, includeDir ...bool) ([]string, error) {
	if !IsDir(rootPath) {
		return nil, errors.New("not a directory or does not exist: " + rootPath)
	}

	isIncludeDir := false
	if len(includeDir) >= 1 {
		isIncludeDir = includeDir[0]
	}
	return statDir(rootPath, "", isIncludeDir, false, true)
}

// GetAllSubDirs returns all subdirectories of given root path.
// Slice does not include given path itself.
func GetAllSubDirs(rootPath string) ([]string, error) {
	if !IsDir(rootPath) {
		return nil, errors.New("not a directory or does not exist: " + rootPath)
	}
	return statDir(rootPath, "", true, true, false)
}

func GetAllFiles(rootPath string) ([]string, error) {
	if !IsDir(rootPath) {
		return nil, errors.New("not a directory or does not exist: " + rootPath)
	}
	return statDir(rootPath, "", true, false, false)
}

// LgetAllSubDirs returns all subdirectories of given root path, including
// following symbolic links, if any.
// Slice does not include given path itself.
func LgetAllSubDirs(rootPath string) ([]string, error) {
	if !IsDir(rootPath) {
		return nil, errors.New("not a directory or does not exist: " + rootPath)
	}
	return statDir(rootPath, "", true, true, true)
}

// GetFileListBySuffix returns an ordered list of file paths.
// It recognize if given path is a file, and don't do recursive find.
func GetFileListBySuffix(dirPath, suffix string) ([]string, error) {
	if !IsExist(dirPath) {
		return nil, fmt.Errorf("given path does not exist: %s", dirPath)
	} else if IsFile(dirPath) {
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

	files := make([]string, 0, len(fis))
	for _, fi := range fis {
		if strings.HasSuffix(fi.Name(), suffix) {
			files = append(files, path.Join(dirPath, fi.Name()))
		}
	}

	return files, nil
}

// CopyDir copy files recursively from source to target directory.
//
// The filter accepts a function that process the path info.
// and should return true for need to filter.
//
// It returns error when error occurs in underlying functions.
func CopyDir(srcPath, destPath string, filters ...func(filePath string) bool) error {
	// Check if target directory exists.
	if IsExist(destPath) {
		return errors.New("file or directory alreay exists: " + destPath)
	}

	err := os.MkdirAll(destPath, os.ModePerm)
	if err != nil {
		return err
	}

	// Gather directory info.
	infos, err := StatDir(srcPath, true)
	if err != nil {
		return err
	}

	var filter func(filePath string) bool
	if len(filters) > 0 {
		filter = filters[0]
	}

	for _, info := range infos {
		if filter != nil && filter(info) {
			continue
		}

		curPath := path.Join(destPath, info)
		if strings.HasSuffix(info, "/") {
			err = os.MkdirAll(curPath, os.ModePerm)
		} else {
			err = CopyFile(path.Join(srcPath, info), curPath)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func MkDir(destPath string) error {
	return os.MkdirAll(destPath, os.ModePerm)
}

func RmDir(destPath string) error {
	return os.RemoveAll(destPath)
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		return ""
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}

// 获取当前执行程序所在的目录名称（如在/home/latelee/foo，返回foo）
func GetRunningDirectory() string {
	exePath, err := os.Executable()
	if err != nil {
		return ""
	}

	ret, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	ret = filepath.Base(ret) // 相对目录
	return ret
}

// 获取文件名称，返回完整文件名、去后缀的部分、后缀(后缀有点号)
// /foo/bar/hello.go 返回：hello.go hello .go
func GetPathFile(dir string) (file, basefile, ext string) {
	file = filepath.Base(dir)
	ext = filepath.Ext(dir)
	basefile = strings.TrimSuffix(file, ext)
	return
}

///////////////////////////////
/*
fileType 0 只有文件
		 1 只有目录
		 2 所有
isInclude 递归
*/
func getAllFiles(dirPath string, fileType int, isInclude bool) (files []string, err error) {
	// fis, err := ioutil.ReadDir(filepath.Clean(filepath.ToSlash(rootPath)))
	// if err != nil {
	// 	return nil, err
	// }
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	fis, err := dir.Readdir(0)
	if err != nil {
		return nil, err
	}

	for _, f := range fis {
		_path := filepath.Join(dirPath, f.Name())

		if f.IsDir() {
			if isInclude {
				fs, _ := getAllFiles(_path, fileType, isInclude)
				files = append(files, fs...)
			}
			if fileType == 1 || fileType == 2 {
				files = append(files, _path)
			}
			continue
		}
		if fileType == 0 || fileType == 2 {
			files = append(files, _path)
		}
	}

	return files, nil
}

func GetAllDirsOneDir(rootPath string) ([]string, error) {
	if !IsDir(rootPath) {
		return nil, errors.New("not a directory or does not exist: " + rootPath)
	}
	return getAllFiles(rootPath, 1, false)
}

func GetAllFilesOneDir(rootPath string) ([]string, error) {
	if !IsDir(rootPath) {
		return nil, errors.New("not a directory or does not exist: " + rootPath)
	}
	return getAllFiles(rootPath, 0, false)
}

func GetAllInDir(rootPath string, include bool) ([]string, error) {
	if !IsDir(rootPath) {
		return nil, errors.New("not a directory or does not exist: " + rootPath)
	}
	return getAllFiles(rootPath, 2, include)
}
