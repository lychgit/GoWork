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

func fileUpload(path string) bool {
	return true;
}