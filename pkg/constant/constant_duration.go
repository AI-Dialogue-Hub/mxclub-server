package constant

import "time"

const (
	Duration_Day   = time.Duration(1) * time.Hour * 24
	Duration_7_Day = Duration_Day * 7
	Duration_5_Day = Duration_Day * 5
	Duration_3_Day = Duration_Day * 3
	Duration_2_Day = Duration_Day * 2
	Duration_1_Day = Duration_Day * 1

	// ============

	Duration_minute    = time.Duration(1) * time.Minute
	Duration_10_minute = Duration_minute * 10
	Duration_20_minute = Duration_minute * 20
)
