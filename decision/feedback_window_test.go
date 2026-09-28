package decision

import (
	"testing"
	"time"
)

func testFeedbackObservation(t *testing.T, kind FeedbackKind, value float64) FeedbackObservation {
	t.Helper()
	observation := FeedbackObservation{
		ReceiptDigest:  "sha256:receipt",
		LedgerDigest:   "sha256:ledger",
		Kind:           kind,
		MetricName:     "support.route.quality",
		MetricValue:    value,
		EvidenceDigest: "sha256:evidence",
		RecordedAt:     time.Date(2026, 9, 28, 0, 0, 0, 0, time.UTC),
		NonAuthorizing: true,
	}
	digest, err := Digest(observation)
	if err != nil {
		t.Fatalf("Digest() error = %v", err)
	}
	observation.FeedbackDigest = digest
	return observation
}

func TestMeasureFeedbackWindowExcludesUnknownFromKnownMean(t *testing.T) {
	window, err := MeasureFeedbackWindow([]FeedbackObservation{
		testFeedbackObservation(t, FeedbackConfirmed, 0.9),
		testFeedbackObservation(t, FeedbackRefuted, 0.4),
		testFeedbackObservation(t, FeedbackUnknown, 0.99),
	}, "support.route.quality")
	if err != nil {
		t.Fatalf("MeasureFeedbackWindow() error = %v", err)
	}
	if window.KnownCount != 2 || window.UnknownCount != 1 {
		t.Fatalf("window counts = %#v, want known=2 unknown=1", window)
	}
	if window.KnownMean != 0.65 {
		t.Fatalf("KnownMean = %v, want 0.65", window.KnownMean)
	}
	if window.ObservedCoverage != 2.0/3.0 {
		t.Fatalf("ObservedCoverage = %v, want 2/3", window.ObservedCoverage)
	}
	if err := window.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}
