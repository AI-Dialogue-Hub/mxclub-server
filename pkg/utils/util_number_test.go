// Copyright 2023 QINIU. All rights reserved
// @Description:
// @Version: 1.0.0
// @Date: 2023/08/16 16:04
// @Author: liangfengyuan@qiniu.com

package utils

import (
	"github.com/stretchr/testify/assert"
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

func TestIsDigit(t *testing.T) {
	assert.True(t, IsNumber("20"))
	assert.True(t, IsNumber("20.5"))
	assert.Equal(t, float64(20), ParseFloat64("20"))
	assert.Equal(t, 20.5, ParseFloat64("20.5"))
}

func TestRoundToDecimalPlaces(t *testing.T) {
	tests := []struct {
		num      float64
		decimals int
		expected float64
	}{
		{123.456789, 2, 123.46},
		{123.456789, 3, 123.457},
		{123.456789, 1, 123.5},
		{250.39999999999998, 2, 250.4},
		{5275.13, 2, 5275.13},
	}

	for _, test := range tests {
		result := RoundToDecimalPlaces(test.num, test.decimals)
		if result != test.expected {
			t.Errorf("For %.2f with %d decimals, expected %.2f, but got %.2f", test.num, test.decimals, test.expected, result)
		}
	}
}

func TestGetDayStartAndEndTimes(t *testing.T) {
	start, end := GetTodayStartAndEndTimes()
	t.Logf("%+v, %+v", start, end)
}
