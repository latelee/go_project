package utils

import (
	//"database/sql"
	//"flag"
	"fmt"
    //"log"
    //"errors"
    //"time"
    //"reflect"
    //"math"
    //"strconv"

    //"strings"
    //"encoding/binary"
    //"io/ioutil"

	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	
)

// TOCHECK：是否要遍历所有字符？
func isGBK(data []byte) bool {
	length := len(data)
    var i int = 0
    for i < length {
        //fmt.Printf("for %x\n", data[i])
        if data[i] <= 127 {
            //编码小于等于127,只有一个字节的编码，兼容ASCII吗
            i++
            continue
        } else {
			if (i+1 >= length) { // 如传入utf8，长度为奇数，下面会越界
				return false;
			}
            //大于127的使用双字节编码
            if  data[i] >= 0x81 &&
                data[i] <= 0xfe &&
                data[i + 1] >= 0x40 &&
                data[i + 1] <= 0xfe &&
                data[i + 1] != 0xf7 {
				i += 2
                continue
            } else {
                return false
            }
        }
    }
    return true
}

func preNUm(data byte) int {
    str := fmt.Sprintf("%b", data)
    var i int = 0
    for i < len(str) {
        if str[i] != '1' {
            break
        }
        i++
    }
    return i
}
func isUtf8(data []byte) bool {
    for i := 0; i < len(data);  {
        if data[i] & 0x80 == 0x00 {
            // 0XXX_XXXX
            i++
            continue
        } else if num := preNUm(data[i]); num > 2 {
            // 110X_XXXX 10XX_XXXX
            // 1110_XXXX 10XX_XXXX 10XX_XXXX
            // 1111_0XXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
            // 1111_10XX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
            // 1111_110X 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
            // preNUm() 返回首个字节的8个bits中首个0bit前面1bit的个数，该数量也是该字符所使用的字节数
            i++
            for j := 0; j < num - 1; j++ {
                //判断后面的 num - 1 个字节是不是都是10开头
                if data[i] & 0xc0 != 0x80 {
                    return false
                }
                i++
            }
        } else  {
            //其他情况说明不是utf-8
            return false
        }
    }
    return true
}

// Golang默认是utf8格式，有些字符串是gbk编码，要转换，否则识别不了
// 注：如果判断已经是utf了？
func GbkToUtf81(s []byte) ([]byte, error) {
	if (isUtf8(s)) {
		return s, nil
	}
    reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
    d, e := ioutil.ReadAll(reader)
    if e != nil {
        return nil, e
    }
    return d, nil
}


func Utf8ToGbk1(s []byte) ([]byte, error) {
	if (isGBK(s)) {
		return s, nil
	}
    reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
    d, e := ioutil.ReadAll(reader)
    if e != nil {
        return nil, e
    }
    return d, nil
}

// 字符串格式
func GbkToUtf8(s1 string) (string, error) {
	s := []byte(s1)
	if (isUtf8(s)) {
		return s1, nil
	}
    reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
    d, e := ioutil.ReadAll(reader)
    if e != nil {
        return "", e
    }

    return string(d[:]), nil
}

func Utf8ToGbk(s1 string) (string, error) {
	s := []byte(s1)
	if (isGBK(s)) {
		return s1, nil
	}
    reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
    d, e := ioutil.ReadAll(reader)
    if e != nil {
        return "", e
    }
    return string(d[:]), nil
}
