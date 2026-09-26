package decision

import (
	"testing"
	"time"
)

func TestFreshAtRejectsStaleAndFutureResults(t *testing.T) {
	result := Result{ObservedAt: time.Unix(100, 0).UTC()}
	now := time.Unix(110, 0).UTC()
	if !result.FreshAt(now, 10*time.Second) {
		t.Fatal("expected result at the freshness boundary to be fresh")
	}
	if result.FreshAt(now.Add(time.Second), 10*time.Second) {
		t.Fatal("expected result beyond the freshness boundary to be stale")
	}
	if result.FreshAt(time.Unix(99, 0).UTC(), 10*time.Second) {
		t.Fatal("future observations must not be treated as fresh")
	}
}
