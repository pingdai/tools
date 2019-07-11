package timelib

import "time"

var (
	MySQLDatetimeZero = MySQLDatetime(time.Time{})
)

// swagger:strfmt date-time
type MySQLDatetime time.Time

func ParseMySQLDatetimeFromString(s string) (dt MySQLDatetime, err error) {
	var t time.Time
	t, err = time.Parse(TIME_FORMAT, s)
	dt = MySQLDatetime(t)
	return
}

func (dt MySQLDatetime) Time() MySQLDatetime {
	return MySQLDatetime(time.Time(dt).In(CST))
}

func (dt MySQLDatetime) String() string {
	return time.Time(dt).In(CST).Format(TIME_FORMAT)
}

func (dt MySQLDatetime) Format(layout string) string {
	return time.Time(dt).In(CST).Format(layout)
}

func (dt MySQLDatetime) MarshalText() ([]byte, error) {
	if dt.IsZero() {
		return []byte(""), nil
	}
	str := dt.String()
	return []byte(str), nil
}

func (dt *MySQLDatetime) UnmarshalText(data []byte) (err error) {
	str := string(data)
	if len(str) == 0 || str == "0" {
		str = MySQLDatetimeZero.String()
	}
	*dt, err = ParseMySQLDatetimeFromString(str)
	return
}

func (dt MySQLDatetime) IsZero() bool {
	unix := dt.Unix()
	return unix == 0 || unix == MySQLDatetimeZero.Unix()
}

func (dt MySQLDatetime) Unix() int64 {
	return time.Time(dt).Unix()
}

func (dt MySQLDatetime) In(loc *time.Location) MySQLDatetime {
	return MySQLDatetime(time.Time(dt).In(loc))
}

// 获取当天最后一秒（东8区）
func (dt MySQLDatetime) GetTodayLastSecCST() MySQLDatetime {
	return MySQLDatetime(GetTodayLastSecInLocation(time.Time(dt), CST))
}

// 添加N个工作日（东8区）
func (dt MySQLDatetime) AddWorkingDaysCST(days int) MySQLDatetime {
	return MySQLDatetime(AddWorkingDaysInLocation(time.Time(dt), days, CST))
}

// 获取当天0点（东8区）
func (dt MySQLDatetime) GetTodayFirstSecCST() MySQLDatetime {
	return MySQLDatetime(GetTodayFirstSecInLocation(time.Time(dt), CST))
}
