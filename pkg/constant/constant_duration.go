package constant

import "time"

const (
	Duration_Day   = time.Duration(1) * time.Hour * 24
	Duration_7_Day = Duration_Day * 7
	Duration_5_Day = Duration_Day * 5
	Duration_3_Day = Duration_Day * 3
)
