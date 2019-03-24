package main

import (
	"strings"
	"time"
	"os"
	"github.com/astaxie/beego"
	"goFrame/utils"
	"encoding/json"
		"bufio"
)

type TempInfo struct {
	ChunkSize int64  `json:"ChunkSize"`
	FileHash  string `json:"FileHash"`
	FileId    string `json:"FileId"`
	FileName  string `json:"FileName"`
	FileSize  int64  `json:"FileSize"`
	Label     string `json:"Label"`
}
const (
	ANSIC       = "Mon Jan _2 15:04:05 2006"
	UnixDate    = "Mon Jan _2 15:04:05 MST 2006"
	RubyDate    = "Mon Jan 02 15:04:05 -0700 2006"
	RFC822      = "02 Jan 06 15:04 MST"
	RFC822Z     = "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
	RFC850      = "Monday, 02-Jan-06 15:04:05 MST"
	RFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
	RFC1123Z    = "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
	RFC3339     = "2006-01-02T15:04:05Z07:00"
	RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
	Kitchen     = "3:04PM"
	// Handy time stamps.
	Stamp      = "Jan _2 15:04:05"
	StampMilli = "Jan _2 15:04:05.000"
	StampMicro = "Jan _2 15:04:05.000000"
	StampNano  = "Jan _2 15:04:05.000000000"

	Year     = "06"
	LongYear = "2006"

	Month     = "Jan"
	ZeroMonth = "01"
	NumMonth  = "1"
	LongMonth = "January"

	Day         = "2"
	ZeroDay     = "02"
	UnderDay    = "_2"
	WeekDay     = "Mon"
	LongWeekDay = "Monday"

	Hour       = "15"
	ZeroHour12 = "03"
	Hour12     = "3"

	Minute     = "4"
	ZeroMinute = "04"

	Second     = "5"
	ZeroSecond = "05"

	PM                    = "PM"
	pm                    = "pm"
	TZ                    = "MST"
	ISO8601TZ             = "Z0700" // prints Z for UTC
	ISO8601SecondsTZ      = "Z070000"
	ISO8601ShortTZ        = "Z07"
	ISO8601ColonTZ        = "Z07:00" // prints Z for UTC
	ISO8601ColonSecondsTZ = "Z07:00:00"
	NumTZ                 = "-0700" // always numeric
	NumSecondsTz          = "-070000"
	NumShortTZ            = "-07"    // always numeric
	NumColonTZ            = "-07:00" // always numeric
	NumColonSecondsTZ     = "-07:00:00"
	FracSecond0           = ".0"               //".00", ... , trailing zeros included
	FracSecond9           = ".9"               //".99", ..., trailing zeros omitted
	stdNeedDate           = 1 << 8             // need month, day, year
	stdNeedClock          = 2 << 8             // need hour, minute, second
	stdArgShift           = 16                 // extra argument in high bits, above low stdArgShift
	stdMask               = 1<<stdArgShift - 1 // mask out argument
)

//根据传入的Y-m-d、 Y/m/d日期格式  生成Go语言time中对应的 显示格式字符串
func format(layout string) string {
	r := strings.NewReplacer("Y", LongYear, "y", Year, "m", ZeroMonth, "d", ZeroDay, "H", Hour, "h", ZeroHour12, "i", ZeroSecond, "s", ZeroMinute)
	return r.Replace(layout)
}

/**
 * t 要格式化的时间
 * layout 要格式化的时间格式
 */
func Date(t time.Time, layout string) string {
	f := format(layout)
	o := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local)
	return o.Format(f)
}

func testFile() {
	taskId := "45a31f103639d71aff37831afd304af3"
	tmpPath := "tests/tmp/"
	infoPath := tmpPath + taskId + "info"
	if utils.IsFile(infoPath) {
		beego.Debug("is a file")
	} else {
		beego.Debug("isn't a file")
	}
	tempInfoFile, err := os.Open(infoPath)
	if err != nil {
		beego.Debug(err.Error())
	}
	var tempInfoSize int64
	if fileInfo, err := os.Stat(infoPath); err == nil {
		tempInfoSize = fileInfo.Size()
	} else {
		beego.Error(err.Error())
	}
	data := TempInfo{}
	var info = make([]byte, tempInfoSize)
	if _, err := tempInfoFile.Read(info); err == nil {
		error := json.Unmarshal(info, &data)
		beego.Error(info)
		if error != nil {
			beego.Error(error.Error())
		}
		beego.Error(data)
	} else {
		beego.Debug(err.Error())
	}

	uploadPath := ""
	mergeFile := "1"
	ext := ".jpg"
	saveDir := utils.Date("Ymd", time.Now()) + string(os.PathSeparator)
	//创建、打开上传文件
	uploadFilePath := uploadPath + saveDir + utils.String(mergeFile) + ext
	uploadFile, err := os.OpenFile(uploadFilePath, os.O_CREATE|os.O_WRONLY, 0777)
	defer uploadFile.Close()
	if err != nil {
		beego.Error("mergeBlock: " + err.Error())
	}
	//锁住文件后合并缓存文件
	i := 0 //文件个数
	var size int64 = 0
	fileSize := data.FileSize
	//优化 -> 加文件锁  ????

	//合并缓存文
	ioW := bufio.NewWriter(uploadFile) //创建新的 Writer 对象
	for size < fileSize {
		chunkFile := tmpPath + taskId  + "." + utils.String(i) + ".tmp"
		tempFile, err := os.OpenFile(chunkFile, os.O_RDONLY, 0777)
		defer func() {
			if tempFile != nil {
				//缓存文件合并过程失败 退出前关闭文件
				tempFile.Close()
			}
		}()
		if err != nil {
			beego.Error("mergeBlock: " + err.Error())
			break
		}
		//将缓存文件内容读取后写入上传文件中
		var buff [] byte
		if _, err := tempFile.Read(buff); err != nil {
			beego.Error("mergeBlock: " + err.Error())
			break
		}
		chunkSize, err := ioW.Write(buff) //缓存块大小
		if err != nil {
			beego.Error("mergeBlock: " + err.Error())
			break
		}
		ioW.Flush()
		size += int64(chunkSize)
		i++
		//合并成功, 关闭缓存文件
		if err := tempFile.Close(); err != nil {
			beego.Error("mergeBlock: " + err.Error())
			break
		}
		//删除合并的缓存文件
		if err := os.Remove(chunkFile); err != nil {
			beego.Error("mergeBlock: " + err.Error())
			break
		}
	}
	if size < fileSize {
		beego.Error("mergeBlock: " + err.Error())
	}
}


func main() {
	//beego.Debug(Date(time.Now(),"Ymd"))
	//
	//data := map[string]string {"filename": "1.jpg"}
	//beego.Error(data["filename"])
	//index := strings.Index(data["filename"], ".")
	//ext := utils.String(data["filename"][index:])
	//beego.Error(ext)
	testFile()

}
