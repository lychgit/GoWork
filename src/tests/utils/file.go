package utils

import (
	"os"
	"fmt"
	"io"
	"crypto/rand"
	"encoding/base64"
		"crypto/md5"
	)

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

//根据传入路径判断文件或目录是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	//file, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
	}
	//beego.Debug(file.IsDir())
	return true
}

//根据传入路径判断文件是否是目录
func IsDir(path string) bool {
	file, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) && file.IsDir() {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
	}
	return file.IsDir()
}

//根据传入路径判断是否是文件
func IsFile(path string) bool {
	file, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) && !file.IsDir() {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
	}
	return !file.IsDir()
}

//根据传入路径判断文件是否可写
func IsWritable(path string) bool {
	return true;
}

//获取路径对应的文件
func GetFile(path string) os.FileInfo {
	file, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) && !file.IsDir() {
			return file
		}
		if os.IsNotExist(err) {
			return nil
		}
	}
	if !file.IsDir() {
		return file
	} else {
		return nil
	}
}
