package util

import "time"

//TimestampFormat Returns　当前时间字符串
//@Summary get now time and format
//@Return (string) YYYY-MM-DD hh:ii:ss
func TimestampFormat() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

//ToTimeFormat Returns 时间转换为字符串
func ToTimeFormat(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func ToTimeFormatPoint(t *time.Time) string {
	if t == nil {
		return "0000-00-00 00:00:00"
	}

	return t.Format("2006-01-02 15:04:05")
}

//Timestamp Returns 当前时间毫秒
func Timestamp() int64 {
	return time.Now().UnixNano() / 1e6
}

//ToTimestamp Returns 返回这个时间的毫秒数
func ToTimestamp(t time.Time) int64 {
	return t.UnixNano() / 1e6
}

//ToTimestampPoint Returns 返回这个时间的毫秒数
func ToTimestampPoint(t *time.Time) int64 {
	if t == nil {
		return 0
	}
	return t.UnixNano() / 1e6
}

//TimestampDay Returns 当前时间日时间间戳,用户潘盾是否是同一天
func TimestampDay() int64 {
	return ToDay(time.Now())
}

//ToDate 字符串时间转换到 Time Date
func ToDate(tm string) time.Time {
	tme, _ := time.Parse("2006-01-02", tm)
	return tme
}

//ToDateTime 字符串转换到 Time DateTime
func ToDateTime(tm string) time.Time {
	tme, _ := time.Parse("2006-01-02 15:04:05", tm)
	return tme
}

//ToDay 转换到时间戳（日）
func ToDay(tm time.Time) int64 {
	return tm.Unix() / 60 / 60 / 24
}

//IsDay 判断两时间段是否是同一天
func IsDay(a time.Time, b time.Time) bool {
	if ToDay(a) == ToDay(b) {
		return true
	}
	return false
}
