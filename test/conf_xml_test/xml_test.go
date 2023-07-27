/*
go test -v -run TestXml

输出：
test of xml...
id: 000250 location: 梧州市岑溪市
filename: /tmp/log.txt time: 0

*/

package test

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"testing"
)

var (
	cfgFile string = "config.xml"
)

type XMLConfig_t struct {
	XMLName    xml.Name    `xml:"config"` // 最外层的标签 config
	InfoNode   XMLInfo_t   `xml:"info"`   // 读取 info 标签下的内容
	ServerNode XMLServer_t `xml:"server"` // 读取 server 标签下的内容
}

type XMLInfo_t struct {
	Id         string      `xml:"id"`
	Num        string      `xml:"num"`
	PersonNode XMLPerson_t `xml:"person"` // 读取 person 标签下的内容
}

type XMLPerson_t struct {
	Location string `xml:"location"`
	Email    string `xml:"email"`
}

type XMLServer_t struct {
	Filename string `xml:"filename"`
	Url      string `xml:"url"`
	Empty    string `xml:"empty"`
	Time     int    `xml:"time"`
}

func TestXml(t *testing.T) {
	fmt.Println("test of xml...")
	// 读文件
	xdata, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		fmt.Printf("open config file [%v] failed: %v\n", cfgFile, err.Error())
		return
	}
	// 解析
	var xmlConfig XMLConfig_t
	err = xml.Unmarshal(xdata, &xmlConfig)
	if err != nil {
		fmt.Printf("parse xml error: %v", err.Error())
		return
	}

	fmt.Printf("id: %v location: %v\n", xmlConfig.InfoNode.Id, xmlConfig.InfoNode.PersonNode.Location)

	fmt.Printf("filename: %v time: %v\n", xmlConfig.ServerNode.Filename, xmlConfig.ServerNode.Time)

}
