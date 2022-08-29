package tcpp

import (
	//"strconv"
	"webdemo/pkg/com"
    //
    //"webdemo/pkg/klog"

)

func checksum(buffer []byte, len int) (sum uint8) {
    var chksum uint8

    for i := 0; i < len; i++ {
        chksum = chksum + buffer[i]
    }
    return chksum&0xff;
}

// 此处先不判断设备ID
func parse_header(buffer []byte, len int) (result int, cmd int, datalen int){

    if len < 5 {
        return -1, -1, -1
    }
    
    if buffer[0] != 0x4c && buffer[0] != 0x45 {
        
        return -2, -1, -1
    }
    
    // 直接转
    datalen = (int)(buffer[4] << 8) + (int)(buffer[3])
    totallen := 5 + datalen

    chksum := checksum(buffer, totallen)
    if chksum != buffer[totallen] {
        return -3, -1, -1
    }
    
    result = 0
    
    cmd = (int)(buffer[2])
    
    return result, cmd, datalen
}

func gen_heartbeat_resp(id string) ([]byte, int) {
    writer := com.NewBufferWriter(100)
    writer.WriteUint8(0x4c)
    writer.WriteUint8(0x45)
    writer.WriteUint8(0x01+0x80)
    writer.WriteUint16(16)
    
    writer.WriteBufferLen([]byte(id), 16)
    
    chksum := checksum(writer.Buffer, writer.Length)
    
    writer.WriteUint8(chksum)
    
    return writer.Buffer, writer.Length
    
}
