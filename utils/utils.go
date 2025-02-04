package utils

import (
	"crypto/md5"
	"encoding/hex"
	"time"
)

func Now() int {
	return int(time.Now().Unix())
}

func EncodeHexMD5(params string) string {
	sumString := md5.Sum([]byte(params))
	return hex.EncodeToString(sumString[:])
}
