package utils

import (
	"time"
	"fmt"
	"strings"
)

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
func Date(layout string, t time.Time) string {
	f := format(layout)
	o := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local)
	return o.Format(f)
}

func GetTime() string {
	const shortForm = "2006-01-01 15:04:05"
	t := time.Now()
	temp := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local)
	str := temp.Format(shortForm)
	fmt.Println(t)
	return str
}

/************/
func TimeParse(t string) int64 {
	p, _ := time.Parse("2006-01-02 15:04:05", t)
	return p.Unix()
}

/***** 时间字符串转换成Time *****/

//不带时区，返回UTC time
func GetTimeFromStr() {
	const format = "2006-01-02 15:04:05"
	timeStr := "2018-01-09 20:24:20"
	p, err := time.Parse(format, timeStr)
	if err == nil {
		fmt.Println(p)
	}
}

//带时区匹配，匹配当前时区的时间
func GetTimeFromStrLoc() {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	const longForm = "Jan 2, 2006 at 3:04pm (MST)"
	t, _ := time.ParseInLocation(longForm, "Jul 9, 2012 at 5:02am (CEST)", loc)
	fmt.Println(t)
	// Note: without explicit zone, returns time in given location.
	const shortForm = "2006-Jan-02"
	t, _ = time.ParseInLocation(shortForm, "2012-Jul-09", loc)
	fmt.Println(t)
}
