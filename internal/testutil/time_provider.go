package testutil

import "time"

type TimeProviderStub struct {
	nowTime time.Time
}

func NewTimeProviderStub(nowTime time.Time) TimeProviderStub {
	return TimeProviderStub{
		nowTime: nowTime,
	}
}

func (t TimeProviderStub) Now() time.Time {
	return t.nowTime
}
