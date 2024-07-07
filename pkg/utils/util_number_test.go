// Copyright 2023 QINIU. All rights reserved
// @Description:
// @Version: 1.0.0
// @Date: 2023/08/16 16:04
// @Author: liangfengyuan@qiniu.com

package utils

import (
	"testing"
	"time"
)

func TestGetCurrentDayStartAndEndDate(t *testing.T) {
	startDate, endDate := GetCurrentDayStartAndEndDate()
	t.Logf("startDate:%v", startDate)
	t.Logf("startDate unix:%v", startDate.Unix())
	t.Logf("endDate:%v", endDate)
	t.Logf("endDate unix:%v", endDate.Unix())
}

func TestConvertTimestampToTime(t *testing.T) {
	// 美国2010/1/1的时间戳
	time, err := ConvertTimestampToTime(int64(1630953600000), "Asia/Shanghai")
	if err != nil {
		t.Logf("err:%v", err.Error())
	}
	t.Logf("time:%v", time)
	time, err = ConvertTimestampToTime(int64(1630953600000), "America/New_York")
	if err != nil {
		t.Logf("err:%v", err.Error())
	}
	t.Logf("time:%v", time)
}

func TestGetYesterdayRange(t *testing.T) {
	start, end := GetYesterdayRange()
	t.Logf("start:%v", start)
	t.Logf("end:%v", end)
}

func TestGetCurrentMonthAndDay(t *testing.T) {
	year, month, day := GetCurrentMonthAndDay()
	t.Logf("year：%d\n", year)
	t.Logf("month：%d\n", month)
	t.Logf("day：%d\n", day)
	t.Logf("%v", 1<<30-1)
}

func TestGetLastMonthStartAndEnd(t *testing.T) {
	t1, t2 := GetLastMonthStartAndEnd()
	t.Logf("%v, %v", t1, t2)
}

func TestParseTimeString(t *testing.T) {
	t.Logf("%v", time.Now().String())
}
