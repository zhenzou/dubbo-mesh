// 重试策略

package retryer

import (
	"math"
	"time"

	"dubbo-mesh/log"
)

type Policy interface {
	Continue() bool
}

func NewTimesPolicy(max int64) Policy {
	return &TimesPolicy{maxAttempts: max, period: 100, maxPeriod: int64(1 * time.Millisecond)}
}

// 限制次数的重试器
type TimesPolicy struct {
	maxAttempts    int64
	attempt        int64
	period         int64
	maxPeriod      int64
	sleptForMillis int64
}

func (t *TimesPolicy) Continue() bool {
	if t.attempt >= t.maxAttempts {
		return false
	}
	log.Debug("retry:", t.attempt)
	t.attempt++
	interval := t.nextMaxInterval()
	time.Sleep(time.Duration(interval) * time.Millisecond)
	t.sleptForMillis += interval
	return true
}

// 重试时间间隔越来越长
func (t *TimesPolicy) nextMaxInterval() int64 {
	interval := int64(float64(t.period) * math.Pow(1.5, float64(t.attempt-1)))
	if interval > t.maxPeriod {
		return t.maxPeriod
	} else {
		return interval
	}
}

func NewDeadlinePolicy(deadline time.Duration) Policy {
	return &DeadlinePolicy{time.Now().Add(deadline)}
}

// 限制时间内
type DeadlinePolicy struct {
	deadline time.Time
}

func (d *DeadlinePolicy) Continue() bool {
	return !time.Now().After(d.deadline)
}
