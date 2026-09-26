package decision

import (
	"testing"
	"time"
)

func TestBuildCounterexampleSetBindsOneAssessment(t *testing.T) {
	confidence := 0.3
	threshold := 0.8
	spec := Spec{
		ID:           "counterexample-set",
		Question:     "Is the proposed change safe enough?",
		Kind:         KindScore,
		Threshold:    &threshold,
		PolicyDigest: "policy-counterexample-set",
	}
	result := Result{
		SpecID:         spec.ID,
		Kind:           spec.Kind,
		Value:          Value{Score: &confidence},
		Confidence:     &confidence,
		EvidenceDigest: "evidence-counterexample-set",
		Provider:       "rule",
		Status:         StatusObserved,
		ObservedAt:     time.Unix(600, 0).UTC(),
	}
	assessment, err := AssessDecision(spec, result, time.Unix(601, 0).UTC(), time.Minute)
	if err != nil {
		t.Fatalf("AssessDecision() error = %v", err)
	}
	first, err := ObserveCounterexample(assessment, "scenario-a", "observation-a", "evidence-a", time.Unix(602, 0).UTC())
	if err != nil {
		t.Fatalf("first counterexample error = %v", err)
	}
	second, err := ObserveCounterexample(assessment, "scenario-b", "observation-b", "evidence-b", time.Unix(603, 0).UTC())
	if err != nil {
		t.Fatalf("second counterexample error = %v", err)
	}
	set, err := BuildCounterexampleSet(assessment, []CounterexampleObservation{first, second})
	if err != nil {
		t.Fatalf("BuildCounterexampleSet() error = %v", err)
	}
	if err := set.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	set.Observations[1].ScenarioDigest = "tampered"
	if err := set.Validate(); err == nil {
		t.Fatal("Validate() accepted a tampered counterexample set")
	}
}
