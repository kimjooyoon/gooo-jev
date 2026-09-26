package decision

import (
	"testing"
	"time"
)

func TestObserveFreshnessClassifiesObservation(t *testing.T) {
	observedAt := time.Unix(100, 0).UTC()
	result := Result{ObservedAt: observedAt}
	fresh := ObserveFreshness(result, time.Unix(105, 0).UTC(), 5*time.Second)
	if fresh.State != FreshnessFresh {
		t.Fatalf("freshness state = %q, want %q", fresh.State, FreshnessFresh)
	}
	stale := ObserveFreshness(result, time.Unix(106, 0).UTC(), 5*time.Second)
	if stale.State != FreshnessStale {
		t.Fatalf("freshness state = %q, want %q", stale.State, FreshnessStale)
	}
	future := ObserveFreshness(result, time.Unix(99, 0).UTC(), 5*time.Second)
	if future.State != FreshnessFuture {
		t.Fatalf("freshness state = %q, want %q", future.State, FreshnessFuture)
	}
	unknown := ObserveFreshness(Result{}, time.Unix(105, 0).UTC(), 5*time.Second)
	if unknown.State != FreshnessUnknown {
		t.Fatalf("freshness state = %q, want %q", unknown.State, FreshnessUnknown)
	}
}
