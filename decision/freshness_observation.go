package decision

import "time"

type FreshnessState string

const (
	FreshnessFresh   FreshnessState = "fresh"
	FreshnessStale   FreshnessState = "stale"
	FreshnessFuture  FreshnessState = "future"
	FreshnessUnknown FreshnessState = "unknown"
)

type FreshnessObservation struct {
	State        FreshnessState
	ObservedAt   time.Time
	EvaluatedAt  time.Time
	MaxAge       time.Duration
}

func ObserveFreshness(result Result, now time.Time, maxAge time.Duration) FreshnessObservation {
	observation := FreshnessObservation{
		State:       FreshnessUnknown,
		ObservedAt:  result.ObservedAt.UTC(),
		EvaluatedAt: now.UTC(),
		MaxAge:      maxAge,
	}
	if result.ObservedAt.IsZero() || now.IsZero() || maxAge < 0 {
		return observation
	}
	if now.Before(result.ObservedAt) {
		observation.State = FreshnessFuture
		return observation
	}
	if result.FreshAt(now, maxAge) {
		observation.State = FreshnessFresh
		return observation
	}
	observation.State = FreshnessStale
	return observation
}
