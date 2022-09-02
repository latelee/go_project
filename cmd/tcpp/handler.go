package tcpp

import (
	//"strconv"
	"webdemo/pkg/com"

	"webdemo/pkg/klog"
)

// TODO：解决粘包问题
func handle(buffer []byte, len int) (obuffer []byte, olen int) {

	dummy := make([]byte, 1)

	ret, cmd, datalen := parse_header(buffer, len)
	if ret < 0 {
		klog.Info("parse errer: ", ret)
		return dummy, 0
	}

	klog.Infof("cmd: %v datalen: %v\n", cmd, datalen)

	reader := com.NewBufferReaderWithLen(buffer[len-datalen-1:], datalen)
	reader.Dump()

	id := reader.ReadString(16)

	switch cmd {
	case 0x01: // 心跳包，回应相同的
		klog.Info("got heartbeat, send back")
		return gen_heartbeat_resp(id)
	case 0x02:
		cameraptz := reader.ReadUint8()
		armleft := reader.ReadUint8()
		armright := reader.ReadUint8()
		temp := reader.ReadUint32()
		hum := reader.ReadUint32()

		klog.Infof("ptz: %v left: %v right: %v temp: %v hum: %v\n",
			cameraptz, armleft, armright, temp, hum)

	}

	return dummy, 0
}
