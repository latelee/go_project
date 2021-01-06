package utils

import (
    uuid "github.com/gofrs/uuid"
    "encoding/hex"
    "crypto/md5"
)

// TO REMOVE...
func Md5(str string) string {
	return Md5ForByte([]byte(str))
}

func Md5ForByte(data []byte) string {
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func GetUUID() string {
	Id := uuid.Must(uuid.NewV4())
	return Id.String()
}

func GenerateKey(userid string) (appKey, appSecret string) {
    appKey = GetUUID()
    tmp := userid + appKey
    appSecret = Md5(tmp)

    return
}