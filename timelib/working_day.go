package timelib

import (
	"math"
	"time"
)

const (
	WORKING_DAYS = 5
	WEEKEND_DAYS = 2
	DAY_SECONDS  = 86400
)

// 根据指定时区获取当天第一秒
func GetTodayFirstSecInLocation(t time.Time, loc *time.Location) time.Time {
	locTime := t.In(loc)
	locTodayMax := time.Date(locTime.Year(), locTime.Month(), locTime.Day(), 0, 0, 0, 0, loc)
	return locTodayMax
}

// 根据指定时区获取当天最后一秒
func GetTodayLastSecInLocation(t time.Time, loc *time.Location) time.Time {
	locTime := t.In(loc)
	locTodayMax := time.Date(locTime.Year(), locTime.Month(), locTime.Day(), 23, 59, 59, 0, loc)
	return locTodayMax
}

// 根据指定时区添加 N 个工作日
func AddWorkingDaysInLocation(t time.Time, days int, loc *time.Location) time.Time {
	locTime := t.In(loc)
	weekDay := locTime.Weekday()
	weekends := int(math.Ceil(float64(int(weekDay)+days)/float64(WORKING_DAYS))) - 1
	weekendDays := weekends * WEEKEND_DAYS
	if weekDay == time.Saturday {
		weekendDays -= 1
	}
	realDays := days + weekendDays
	return locTime.AddDate(0, 0, realDays)
}

// 计算日期偏差
func CountDateDiff(startTime, endTime time.Time, loc *time.Location) int64 {
	return (GetTodayFirstSecInLocation(endTime, loc).Unix() - GetTodayFirstSecInLocation(startTime, loc).Unix()) / DAY_SECONDS
}
