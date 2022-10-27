package time

//直接取什么格式的时间
//直接拿一个星期、月、年后的时间
//其他一些快捷函数

import (
	"fmt"
	"time"
)

//Now 当前时间戳
func Now() int64 {
	return time.Now().Unix()
}

//NowTime 当前时分秒
func NowTime() string {
	return time.Now().Format("15:04:05")
}

//NowDate 当前日期
func NowDate() string {
	return time.Now().Format("2006-01-02")
}

//NowDateTime 当前日期时间
func NowDateTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

//ToTime 时间戳转成时分秒
func ToTime(ti int64) string {
	return time.Unix(ti, 0).Format("15:04:05")
}

//ToDate 时间戳转成日期
func ToDate(ti int64) string {
	return time.Unix(ti, 0).Format("2006-01-02")
}

//ToDateTime 时间戳转成日期时间
func ToDateTime(ti int64) string {
	return time.Unix(ti, 0).Format("2006-01-02 15:04:05")
}

//T0Time 今天零点时间戳
func T0Time() int64 {
	currentTime := time.Now()
	startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
	return startTime.Unix()
}

//T24Time 今天23：59：59时间戳
func T24Time() int64 {
	currentTime := time.Now()
	startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, currentTime.Location())
	return startTime.Unix()
}

//StrToTime 2006-01-02 15:04:05格式的字符串转成时间戳
func StrToTime(s string) int64 {
	ti, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
	if err == nil {
		return ti.Unix()
	}

	return 0
}

//GetM1Time 获取某年某月1号0点的时间
func GetM1Time(year, month int) (time.Time, error) {
	return time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%04d-%02d-01 00:00:00", year, month), time.Local)
}

//DestTime 在当前时间基础上增加或减少时间后转成目标时间戳 如：2h -3h 50m -5m 36s
func DestTime(s string) int64 {
	currentTime := time.Now()
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0
	}
	result := currentTime.Add(d)
	return result.Unix()
}
