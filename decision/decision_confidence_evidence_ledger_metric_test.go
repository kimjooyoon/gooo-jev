package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerMetricIsOneHot(t *testing.T) {
	entry, err := NewDecisionConfidenceEvidenceLedgerEntry(1, "decision", "source-a", "")
	if err != nil {
		t.Fatalf("create entry: %v", err)
	}
	metric := ObserveDecisionConfidenceEvidenceLedgerMetric(entry)
	if metric.Status != "verified" || metric.VerifiedCount != 1 || metric.BrokenCount != 0 || !metric.NonAuthorizing {
		t.Fatalf("unexpected verified metric: %#v", metric)
	}

	tampered := entry
	tampered.SourceDigest = "source-tampered"
	metric = ObserveDecisionConfidenceEvidenceLedgerMetric(tampered)
	if metric.Status != "unknown" || metric.VerifiedCount != 0 || metric.BrokenCount != 1 || !metric.NonAuthorizing {
		t.Fatalf("unexpected unknown metric: %#v", metric)
	}
}
