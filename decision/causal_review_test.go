package decision

import (
	"testing"
	"time"
)

func causalHistorySummary(t *testing.T, mean float64, at time.Time) FeedbackSummary {
	t.Helper()
	summary := FeedbackSummary{
		MetricName:     "utility.delta",
		Count:          1,
		ConfirmedCount: 1,
		MetricSum:      mean,
		MetricMean:     mean,
		MetricMin:      mean,
		MetricMax:      mean,
		LastRecordedAt: at,
		NonAuthorizing: true,
	}
	var err error
	summary.SummaryDigest, err = Digest(summary)
	if err != nil {
		t.Fatalf("Digest() error = %v", err)
	}
	return summary
}

func TestBuildCausalReviewBindsAllEvidence(t *testing.T) {
	confidence := 0.3
	threshold := 0.8
	spec := Spec{
		ID:           "causal-review",
		Question:     "Is the change safe enough?",
		Kind:         KindScore,
		Threshold:    &threshold,
		PolicyDigest: "policy-causal-review",
	}
	result := Result{
		SpecID:         spec.ID,
		Kind:           spec.Kind,
		Value:          Value{Score: &confidence},
		Confidence:     &confidence,
		EvidenceDigest: "evidence-causal-review",
		Provider:       "rule",
		Status:         StatusObserved,
		ObservedAt:     time.Unix(700, 0).UTC(),
	}
	assessment, err := AssessDecision(spec, result, time.Unix(701, 0).UTC(), time.Minute)
	if err != nil {
		t.Fatalf("AssessDecision() error = %v", err)
	}
	counterexample, err := ObserveCounterexample(assessment, "scenario-causal", "observation-causal", "evidence-causal", time.Unix(702, 0).UTC())
	if err != nil {
		t.Fatalf("ObserveCounterexample() error = %v", err)
	}
	counterexamples, err := BuildCounterexampleSet(assessment, []CounterexampleObservation{counterexample})
	if err != nil {
		t.Fatalf("BuildCounterexampleSet() error = %v", err)
	}
	history, err := BuildFeedbackHistory([]FeedbackSummary{
		causalHistorySummary(t, 0.2, time.Unix(703, 0).UTC()),
		causalHistorySummary(t, 0.1, time.Unix(704, 0).UTC()),
	}, time.Unix(705, 0).UTC())
	if err != nil {
		t.Fatalf("BuildFeedbackHistory() error = %v", err)
	}
	review, err := BuildCausalReview(assessment, counterexamples, history)
	if err != nil {
		t.Fatalf("BuildCausalReview() error = %v", err)
	}
	if err := review.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if review.Status != CausalReviewReview || review.CausalReason != "COUNTEREXAMPLES_PRESENT" {
		t.Fatalf("review = %#v", review)
	}
	review.CounterexampleSetDigest = "tampered"
	if err := review.Validate(); err == nil {
		t.Fatal("Validate() accepted a tampered causal review")
	}
}
