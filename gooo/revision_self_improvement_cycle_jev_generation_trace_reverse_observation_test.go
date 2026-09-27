package gooo

import "testing"

func TestRevisionSelfImprovementCycleJEVGenerationTraceReverseObservationPreservesUnknown(t *testing.T) {
	generationObservation := ObserveRevisionSelfImprovementCycleJEVGenerationTraceObservation(
		RevisionSelfImprovementCycleJEVGenerationTraceEvidence{},
	)
	reverseObservation := ObserveRevisionSelfImprovementCycleJEVGenerationTraceReverseObservation(generationObservation)

	if reverseObservation.Status != "UNKNOWN" {
		t.Fatalf("status = %q, want UNKNOWN", reverseObservation.Status)
	}
	if reverseObservation.MissingStage == "" {
		t.Fatal("missing stage must be preserved")
	}
	if reverseObservation.ReconstructedTraceDigest != "" {
		t.Fatal("unknown reverse observation must not fabricate a reconstruction")
	}
	if err := reverseObservation.Validate(); err != nil {
		t.Fatalf("reverse observation should validate: %v", err)
	}
}

func TestRevisionSelfImprovementCycleJEVGenerationTraceReverseObservationRejectsUnboundTrace(t *testing.T) {
	reverseObservation := ObserveRevisionSelfImprovementCycleJEVGenerationTraceReverseObservation(
		RevisionSelfImprovementCycleJEVGenerationTraceObservation{Status: "BOUND"},
	)

	if reverseObservation.Status != "UNKNOWN" {
		t.Fatalf("status = %q, want UNKNOWN", reverseObservation.Status)
	}
	if reverseObservation.MissingStage == "" {
		t.Fatal("unbound trace must preserve an unresolved stage")
	}
}