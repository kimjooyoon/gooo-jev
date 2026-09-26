package decision

import (
	"testing"
	"time"
)

func TestObserveCounterexampleBindsAssessmentEvidence(t *testing.T) {
	confidence := 0.4
	spec := Spec{
		ID:           "counterexample-review",
		Question:     "Is the proposed change safe enough?",
		Kind:         KindScore,
		Threshold:    func() *float64 { value := 0.8; return &value }(),
		PolicyDigest: "policy-counterexample",
	}
	result := Result{
		SpecID:         spec.ID,
		Kind:           spec.Kind,
		Value:          Value{Score: &confidence},
		Confidence:     &confidence,
		EvidenceDigest: "evidence-counterexample",
		Provider:       "rule",
		Status:         StatusObserved,
		ObservedAt:     time.Unix(500, 0).UTC(),
	}
	assessment, err := AssessDecision(spec, result, time.Unix(501, 0).UTC(), time.Minute)
	if err != nil {
		t.Fatalf("AssessDecision() error = %v", err)
	}
	counterexample, err := ObserveCounterexample(
		assessment,
		"scenario-counterexample",
		"observation-counterexample",
		"evidence-counterexample-observed",
		time.Unix(502, 0).UTC(),
	)
	if err != nil {
		t.Fatalf("ObserveCounterexample() error = %v", err)
	}
	if err := counterexample.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if !counterexample.NonAuthorizing || counterexample.AssessmentDigest != assessment.AssessmentDigest {
		t.Fatalf("counterexample = %#v", counterexample)
	}
	counterexample.ScenarioDigest = "tampered"
	if err := counterexample.Validate(); err == nil {
		t.Fatal("Validate() accepted a tampered counterexample")
	}
}
