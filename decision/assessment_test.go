package decision

import (
	"testing"
	"time"
)

func TestAssessDecisionPreservesReviewAndUnknownStates(t *testing.T) {
	threshold := 0.8
	spec := Spec{
		ID:           "assessment-review",
		Question:     "Is the proposed change safe enough?",
		Kind:         KindScore,
		Threshold:    &threshold,
		PolicyDigest: "policy-assessment",
	}
	confidence := 0.9
	result := Result{
		SpecID:         spec.ID,
		Kind:           spec.Kind,
		Value:          Value{Score: &confidence},
		Confidence:     &confidence,
		EvidenceDigest: "evidence-assessment",
		Provider:       "rule",
		Status:         StatusObserved,
		ObservedAt:     time.Unix(300, 0).UTC(),
	}
	observed, err := AssessDecision(spec, result, time.Unix(301, 0).UTC(), time.Minute)
	if err != nil {
		t.Fatalf("AssessDecision() error = %v", err)
	}
	if observed.Status != AssessmentObserved {
		t.Fatalf("status = %q, want %q", observed.Status, AssessmentObserved)
	}
	stale, err := AssessDecision(spec, result, time.Unix(400, 0).UTC(), time.Minute)
	if err != nil {
		t.Fatalf("stale assessment error = %v", err)
	}
	if stale.Status != AssessmentReview || stale.CausalReason != "DECISION_FRESHNESS_STALE" {
		t.Fatalf("stale = %#v", stale)
	}
	result.Status = StatusUnknown
	unknown, err := AssessDecision(spec, result, time.Unix(301, 0).UTC(), time.Minute)
	if err != nil {
		t.Fatalf("unknown assessment error = %v", err)
	}
	if unknown.Status != AssessmentUnknown {
		t.Fatalf("unknown = %#v", unknown)
	}
	if err := unknown.Validate(); err != nil {
		t.Fatalf("unknown validation error = %v", err)
	}
}
