package utils

import (
	"os"
	"fmt"
	"io"
	"crypto/rand"
	"encoding/base64"
		"crypto/md5"
)

func isExist(path string)(bool){
	_, err := os.Stat(path)
	if err != nil{
		if os.IsExist(err){
			return true
		}
		if os.IsNotExist(err){
			return false
		}
		fmt.Println(err)
		return false
	}
	return true
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err!= nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func Md5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

//生成Guid字串
func UniqueId() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return Md5(base64.URLEncoding.EncodeToString(b))
}