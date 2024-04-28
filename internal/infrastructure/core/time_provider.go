package core

import "time"

type timeProvider struct{}

func (t timeProvider) Now() time.Time {
	return time.Now()
}
