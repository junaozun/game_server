package global

import (
	"time"
)

type TimeOffsetSeconds int64

func (o *TimeOffsetSeconds) Clear() {
	*o = 0
}

func (o *TimeOffsetSeconds) SetSeconds(seconds int) {
	*o = TimeOffsetSeconds(seconds)
}

func (o *TimeOffsetSeconds) SetHour(h int) {
	*o = TimeOffsetSeconds(h * 3600)
}

func (o *TimeOffsetSeconds) AddHour(h int) {
	*o += TimeOffsetSeconds(h * 3600)
}

func (o *TimeOffsetSeconds) AddMinute(m int) {
	*o += TimeOffsetSeconds(m * 60)
}

func (o *TimeOffsetSeconds) AddSecond(s int) {
	*o += TimeOffsetSeconds(s)
}

func (o *TimeOffsetSeconds) SetDate(t time.Time) {
	*o = TimeOffsetSeconds(t.Unix() - time.Now().Unix())
}

func (o *TimeOffsetSeconds) SetTimeStamp(t int64) {
	*o = TimeOffsetSeconds(t - time.Now().Unix())
}

var TimeOffset TimeOffsetSeconds

// Now 当前时间
func Now() time.Time {
	if TimeOffset == 0 {
		return time.Now()
	}
	return time.Now().Add(time.Second * time.Duration(TimeOffset))
}
