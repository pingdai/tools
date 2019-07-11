package timelib

import "time"

const (
	CST_ZONE_OFFSET = 8 * 60 * 60
	TIME_FORMAT     = "2006-01-02 15:04:05"
)

func Now() time.Time {
	return time.Now().In(CST)
}

var (
	UTC *time.Location
	CST *time.Location
)

func init() {
	UTC = time.UTC
	CST = time.FixedZone("CST", CST_ZONE_OFFSET)
}
