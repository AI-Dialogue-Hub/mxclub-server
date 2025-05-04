// Copyright 2023 QINIU. All rights reserved
// @Description: 注意：使用者请务必确定类型转换会成功！！
// @Version: 1.0.0
// @Date: 2023/05/09 13:15
// @Author: liangfengyuan@qiniu.com

package utils

import (
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"log"
	"math"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

//================  number相关  ====================

func ParseString(args interface{}) string {
	if args == nil {
		log.Panic("NPE")
	}
	// 判断类型（未出现的类型可能会转换失败，请测试）
	switch reflect.TypeOf(args).Kind() {
	case reflect.Float32:
		return strconv.FormatFloat(float64(args.(float32)), 'f', -1, 32)
	case reflect.Float64:
		return strconv.FormatFloat(args.(float64), 'f', -1, 64)
	case reflect.String:
		return args.(string)
	case reflect.Int:
		return strconv.Itoa(args.(int))
	case reflect.Int32:
		return Int32ToStr(args.(int32))
	case reflect.Int64:
		return Int64ToStr(args.(int64))
	case reflect.Uint64:
		return strconv.FormatUint(args.(uint64), 10)
	case reflect.Uint:
		return strconv.FormatUint(uint64(args.(uint)), 10)
	default:
		// 对于未知类型，使用 fmt 包中的 %v 标记输出
		return fmt.Sprintf("%v", args)
	}
}

func Int64ToStr(val int64) string {
	return strconv.FormatInt(val, 10)
}

func Int32ToStr(val int32) string {
	buf := [11]byte{}
	pos := len(buf)
	i := int64(val)
	signed := i < 0
	if signed {
		i = -i
	}
	for {
		pos--
		buf[pos], i = '0'+byte(i%10), i/10
		if i == 0 {
			if signed {
				pos--
				buf[pos] = '-'
			}
			return string(buf[pos:])
		}
	}
}

func IntToStr(val int) string {
	return strconv.Itoa(val)
}

func ParseInt64(args interface{}) int64 {
	return ParseDefaultInt(args, func(str string) (interface{}, error) {
		return strconv.ParseInt(str, 10, 64)
	}).(int64)
}

func ParseInt32(args interface{}) int32 {
	return int32(ParseInt64(args))
}

func ParseUint(args interface{}) uint {
	defer RecoverByPrefixNoCtx("ParseUint")
	return uint(ParseInt64(args))
}

func ParseUint8(args interface{}) uint8 {
	return uint8(ParseInt64(args))
}

func ParseUint32(args interface{}) uint32 {
	return uint32(ParseInt64(args))
}

func ParseUint64(args interface{}) uint64 {
	return uint64(ParseInt64(args))
}

func SafeParseUint64(args interface{}) (result uint64) {
	defer func() {
		if err := recover(); err != nil {
			xlog.Errorf("SafeParseUint64 ERROR:%v", err)
			result = 0
		}
	}()
	result = uint64(ParseInt64(args))
	return
}

func SafeParseNumber[T uint | uint32 | uint64 | int | int32 | int64](args any) (result T) {
	defer func() {
		if err := recover(); err != nil {
			xlog.Errorf("SafeParseUint64 ERROR:%v", err)
			result = 0
		}
	}()
	result = T(ParseInt64(args))
	return
}

func ParseInt(args interface{}) int {
	return ParseDefaultInt(args, func(str string) (interface{}, error) {
		return strconv.Atoi(str)
	}).(int)
}

func ParseFloat64(str string) (result float64) {
	floatVal, err := strconv.ParseFloat(strings.TrimSpace(str), 64)
	if err != nil {
		xlog.Errorf("ParseFloat64 ERROR:%v", err)
		return
	}
	result = floatVal
	return
}

func ParseDefaultInt(args interface{}, callBack func(str string) (interface{}, error)) interface{} {
	if args == nil {
		log.Panic("NPE")
	}
	intArgs, err := callBack(ParseString(args))
	if err != nil {
		log.Panicf("类型转换失败，错误信息：%v", err)
	}
	return intArgs
}

// ParseType
//
//	@Description: 将args转为指定的类型【stringType】
//	@args args 需要类型转换的参数
//	@args stringType 要转换的类型
//	@return newArgs 返回转换好类型的参数
func ParseType(args interface{}, stringType string) (newArgs interface{}) {
	switch stringType {
	case "int":
		return ParseInt(args)
	case "int32":
		return ParseInt32(args)
	case "int64":
		return ParseInt64(args)
	case "uint8":
		return ParseUint8(args)
	case "uint32":
		return ParseUint32(args)
	case "uint64":
		return ParseUint64(args)
	case "string":
		fallthrough
	default:
		return ParseString(args)
	}
}

func IsNotDigit(str string) bool {
	return !IsDigit(str)
}

// IsDigit 检查字符串是否为数字字符串
func IsDigit(str string) bool {
	for _, x := range []rune(str) {
		if !unicode.IsDigit(x) {
			return false
		}
	}
	return true
}

var numberRegex = regexp.MustCompile(`^[-+]?\d*\.?\d+$`)

func IsNumber(s string) bool {
	return numberRegex.MatchString(s)
}

func Max[T int | int32 | int64](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Min[T int | int32 | int64](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// ParseTimeString 将字符串时间格式为`yyyyMMdd`
func ParseTimeString(timeString string) (time.Time, error) {
	return time.Parse("20060102", timeString)
}

// GetLastWeek 返回上一周开始时间和结束时间
func GetLastWeek() (time.Time, time.Time) {
	// 获取当前时间
	now := time.Now()

	// 计算相对于当前时间的上一周开始和结束日期
	weekAgo := now.AddDate(0, 0, -7)
	weekStart := weekAgo.AddDate(0, 0, -int(weekAgo.Weekday())+1)
	weekEnd := weekAgo.AddDate(0, 0, 7-int(weekAgo.Weekday()))

	return weekStart, weekEnd
}

func GetLastMonthStartAndEnd() (time.Time, time.Time) {
	now := time.Now()
	year, month, _ := now.AddDate(0, -1, 0).Date()
	// 加载本地时区
	loc, _ := time.LoadLocation("")
	// 获取上个月的开始时间
	startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, loc)

	// 获取上个月的结束时间
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	return startOfMonth, endOfMonth
}

func GetMonthStartAndEnd(dateStr string) (time.Time, time.Time, error) {
	layout := "2006-01"
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	startOfMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	return startOfMonth, endOfMonth, nil
}

const (
	KindDay   = "day"
	KindWeek  = "week"
	KindMonth = "month"
	KindYear  = "year"
)

func GetDefaultStartAndEndDate(kind string) (startDate, endDate time.Time, err error) {
	// 加载本地时区
	loc, err := time.LoadLocation("")
	if err != nil {
		fmt.Println(err)
		return
	}
	return GetStartAndEndDate(kind, loc.String())
}

// GetStartAndEndDate 获取本周 本月 本年的开始和结束时间
// 注意：由于时区不同，只能保证年月日正确，时分秒可能不正确
func GetStartAndEndDate(kind, timezone string) (startDate, endDate time.Time, err error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("failed to load location: %v", err)
	}

	now := time.Now().In(loc)
	switch kind {
	case KindDay:
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endDate = startDate.Add(24 * time.Hour).Add(-time.Nanosecond)
	case KindWeek:
		startDate = now.AddDate(0, 0, -int(now.Weekday())+1) // 获取本周的开始日期
		endDate = startDate.AddDate(0, 0, 6)                 // 获取本周的结束日期
	case KindMonth:
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc) // 获取本月的开始日期
		endDate = startDate.AddDate(0, 1, 0).Add(-time.Nanosecond)         // 获取本月的结束日期
	case KindYear:
		startDate = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, loc)                     // 获取本年的开始日期
		endDate = time.Date(now.Year(), 12, 31, 23, 59, 59, int(time.Second)-1, loc) // 获取本年的结束日期
	default:
		return time.Time{}, time.Time{}, fmt.Errorf("invalid kind: %s", kind)
	}
	return startDate, endDate, nil
}

// GetCurrentDayStartAndEndDate 本日开始和结束日期
// 注意：由于时区不同，只能保证年月日正确，时分秒可能不正确
func GetCurrentDayStartAndEndDate() (startDate, endDate time.Time) {
	t1, t2, _ := GetDefaultStartAndEndDate(KindDay)
	return t1, t2
}

// GetCurrentWeekStartAndEndDate 本周开始和结束日期
// 注意：由于时区不同，只能保证年月日正确，时分秒可能不正确
func GetCurrentWeekStartAndEndDate() (startDate, endDate time.Time) {
	t1, t2, _ := GetDefaultStartAndEndDate(KindWeek)
	return t1, t2
}

// GetCurrentMonthStartAndEndDate 本月开始和结束日期
// 注意：由于时区不同，只能保证年月日正确，时分秒可能不正确
func GetCurrentMonthStartAndEndDate() (startDate, endDate time.Time) {
	t1, t2, _ := GetDefaultStartAndEndDate(KindMonth)
	return t1, t2
}

// GetCurrentYearStartAndEndDate 本年开始和结束日期
// 注意：由于时区不同，只能保证年月日正确，时分秒可能不正确
func GetCurrentYearStartAndEndDate() (startDate, endDate time.Time) {
	t1, t2, _ := GetDefaultStartAndEndDate(KindYear)
	return t1, t2
}

// GetMonthStartAndEndDate 获取指定月份开始和结束的日期
func GetMonthStartAndEndDate(month int) (time.Time, time.Time) {
	now := time.Now()
	year := now.Year()
	startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	endOfMonth := startOfMonth.AddDate(0, 1, -1)
	return startOfMonth, endOfMonth
}

// GetDaysInMonth 获取这个月有多少天
func GetDaysInMonth() int {
	now := time.Now()
	year, month, _ := now.Date()
	// 获取本月最后一天
	lastDayOfMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.Local)
	// 返回天数
	return lastDayOfMonth.Day()
}

// ConvertTimestampToTime 将传入的时间戳和时区转换为本地的时间
//
//	@args timestamp 毫秒时间戳 有13位
//	@args timezone 时区 例如`Asia/Shanghai`，为空加载本地时区
//	@return time.Time 返回time.Time
//	@return error
func ConvertTimestampToTime(timestamp int64, timezone string) (time.Time, error) {
	if timezone == "" {
		loc, _ := time.LoadLocation("")
		timezone = loc.String()
	}
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}

	// 转换为本地时间
	t := time.Unix(timestamp/1000, (timestamp%1000)*int64(time.Millisecond)).In(loc)

	return t, nil
}

// GetYesterdayRange 获取昨天开始和结束的时间
func GetYesterdayRange() (time.Time, time.Time) {
	// 获取当前时间
	now := time.Now()

	// 计算昨天的日期
	yesterday := now.AddDate(0, 0, -1)

	// 设置昨天的开始时间为当天的 00:00:00
	yesterdayStart := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, now.Location())

	// 设置昨天的结束时间为当天的 23:59:59
	yesterdayEnd := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, 0, now.Location())

	return yesterdayStart, yesterdayEnd
}

func GetCurrentMonthAndDay() (int, int, int) {
	now := time.Now()
	year, month, day := now.Date()
	return year, int(month), day
}

// GetDaysInMonthWithDateStr dateStr => 2023-08
func GetDaysInMonthWithDateStr(dateStr string) int {
	layout := "2006-01"
	t, _ := time.Parse(layout, dateStr)

	// 获取下个月的第一天
	nextMonth := t.AddDate(0, 1, 0)
	// 获取当前月的最后一天
	lastDayOfMonth := nextMonth.AddDate(0, 0, -1)

	return lastDayOfMonth.Day()
}

func ValidateDate(dateStr string) bool {
	// 定义日期字符串的正则表达式，格式为YYYY-MM-DD
	regex := `^\d{4}-\d{2}-\d{2}$`
	match, _ := regexp.MatchString(regex, dateStr)
	return match
}

// GetDayStartAndEndTimes 给定时间返回这天开始和结束的时间
func GetDayStartAndEndTimes(dateStr string) (time.Time, time.Time, error) {
	if !ValidateDate(dateStr) {
		return time.Time{}, time.Time{}, fmt.Errorf("日期字符串不符合格式要求")
	}
	layout := "2006-01-02"
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	startTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	endTime := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())

	return startTime, endTime, nil
}

// GetTodayStartAndEndTimes 给定当天时间返回这天开始和结束的时间
func GetTodayStartAndEndTimes() (time.Time, time.Time) {
	var t = time.Now()

	startTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	endTime := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())

	return startTime, endTime
}

func GetRandomInt(min, max int) int {
	if min > max {
		panic("min cannot more than max")
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func RoundToDecimalPlaces(num float64, decimals int) float64 {
	shift := math.Pow(10, float64(decimals))
	return math.Round(num*shift) / shift
}

func RoundToTwoDecimalPlaces(num float64) float64 {
	return RoundToDecimalPlaces(num, 2)
}
