package decision

import "time"

func (result Result) FreshAt(now time.Time, maxAge time.Duration) bool {
	if result.ObservedAt.IsZero() || now.IsZero() || maxAge < 0 {
		return false
	}
	observedAt := result.ObservedAt.UTC()
	now = now.UTC()
	if now.Before(observedAt) {
		return false
	}
	return now.Sub(observedAt) <= maxAge
}
