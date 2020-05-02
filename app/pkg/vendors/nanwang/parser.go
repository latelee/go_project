package nanwang

import (
    "webdemo/pkg/com"
    "k8s.io/klog"

    //"strings"
)

func parsePacket(bin []byte) {

    klog.Printf("foo")
/*
    var buf mybuffer.BufferReader;
    array := make(map[string]interface{})

    buf.SkipBytes(1);

    klog.Printf("array:\n%##v\n", array);
*/
}

// 测试
func Parser() {

    var str = "681c00682cf6028402010000000001000441c10000000101000441c1000000"
    parsePacket([]byte(com.ToHexByte(str)));
 
}