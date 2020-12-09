package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: MD5V
//@description: md5加密
//@param: str []byte
//@return: string

func MD5V(str []byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}

func Struct2Md5(stc interface{}) string {
	b, _ := json.Marshal(stc)
	h := md5.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}
