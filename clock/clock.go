package clock

import "time"

func RoundToNearestMinute(d time.Duration) time.Duration {
	now := time.Now()
	nextTick := now.Truncate(time.Minute).Add(d)
	return nextTick.Sub(now)
}

func RoundTimeToPreviousMinute(currTime time.Time) time.Time {
	return currTime.Add(-time.Duration(currTime.Second()) * time.Second)
}
