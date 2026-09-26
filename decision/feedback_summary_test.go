package decision

import (
	"math"
	"testing"
	"time"
)

func TestSummarizeFeedbackProducesBoundedMetrics(t *testing.T) {
	spec := Spec{
		ID:             "summary-review",
		Question:       "Did the feedback improve the result?",
		Kind:           KindChoice,
		AllowedChoices: []string{"yes", "no"},
		PolicyDigest:   "policy-summary",
	}
	result := Result{
		SpecID:         spec.ID,
		Kind:           spec.Kind,
		Value:          Value{Choice: "yes"},
		EvidenceDigest: "evidence-summary",
		Provider:       "rule",
		Status:         StatusObserved,
		ObservedAt:     time.Unix(150, 0).UTC(),
	}
	ledger, err := (Ledger{}).AppendVerified(spec, State{Digest: "state-summary"}, result)
	if err != nil {
		t.Fatalf("AppendVerified() error = %v", err)
	}
	first, err := ObserveFeedback(ledger, ledger.Entries[0].Receipt, FeedbackConfirmed, "utility.delta", 0.25, "feedback-1", time.Unix(160, 0).UTC())
	if err != nil {
		t.Fatalf("first feedback error = %v", err)
	}
	second, err := ObserveFeedback(ledger, ledger.Entries[0].Receipt, FeedbackRefuted, "utility.delta", -0.1, "feedback-2", time.Unix(170, 0).UTC())
	if err != nil {
		t.Fatalf("second feedback error = %v", err)
	}
	summary, err := SummarizeFeedback([]FeedbackObservation{first, second}, "utility.delta")
	if err != nil {
		t.Fatalf("SummarizeFeedback() error = %v", err)
	}
	if err := summary.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if summary.Count != 2 || summary.ConfirmedCount != 1 || summary.RefutedCount != 1 ||
		summary.UnknownCount != 0 || math.Abs(summary.MetricMean-0.075) > 1e-9 {
		t.Fatalf("summary = %#v", summary)
	}
	summary.MetricMax = 999
	if err := summary.Validate(); err == nil {
		t.Fatal("Validate() accepted a tampered metric")
	}
}
