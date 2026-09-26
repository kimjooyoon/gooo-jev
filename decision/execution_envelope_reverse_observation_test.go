package decision

import "testing"

func TestObserveExecutionEnvelopeProvenanceReverseReproducesExactBoundary(t *testing.T) {
	result := ObserveExecutionEnvelopeProvenanceReverse(ExecutionEnvelopeReverseObservationInput{
		ObservedStatus:         "ready",
		ExpectedStatus:         "ready",
		ObservedEvidenceDigest: "digest",
		ExpectedEvidenceDigest: "digest",
		NonAuthorizing:         true,
	})
	if result.Status != "reproduced" || result.FirstMismatch != "" || result.EvidenceDigest != "digest" || !result.NonAuthorizing {
		t.Fatalf("complete reverse boundary was not preserved: %#v", result)
	}
}

func TestObserveExecutionEnvelopeProvenanceReverseReportsFirstCounterexample(t *testing.T) {
	result := ObserveExecutionEnvelopeProvenanceReverse(ExecutionEnvelopeReverseObservationInput{
		ObservedStatus:         "ready",
		ExpectedStatus:         "hold",
		ObservedEvidenceDigest: "observed",
		ExpectedEvidenceDigest: "expected",
		NonAuthorizing:         true,
	})
	if result.Status != "counterexample" || result.FirstMismatch != "status" || result.EvidenceDigest != "observed" || !result.NonAuthorizing {
		t.Fatalf("first reverse counterexample was not preserved: %#v", result)
	}
}

func TestObserveExecutionEnvelopeProvenanceReverseRejectsAuthorizationClaim(t *testing.T) {
	result := ObserveExecutionEnvelopeProvenanceReverse(ExecutionEnvelopeReverseObservationInput{
		ObservedStatus: "ready",
		ExpectedStatus: "ready",
		NonAuthorizing: false,
	})
	if result.Status != "UNKNOWN" || result.FirstMismatch != "authorization-boundary" || result.NonAuthorizing {
		t.Fatalf("authorization claim escaped reverse observation: %#v", result)
	}
}

func TestBindExecutionEnvelopeReverseObservationReplaysObservationDigest(t *testing.T) {
	reproduced := ObserveExecutionEnvelopeProvenanceReverse(ExecutionEnvelopeReverseObservationInput{
		ObservedStatus:         "ready",
		ExpectedStatus:         "ready",
		ObservedEvidenceDigest: "digest",
		ExpectedEvidenceDigest: "digest",
		NonAuthorizing:         true,
	})
	binding := BindExecutionEnvelopeReverseObservation(reproduced)
	if binding.Status != "reproduced" || binding.ObservationDigest == "" || binding.ObservedEvidenceDigest != "digest" {
		t.Fatalf("reverse observation was not sealed: %#v", binding)
	}
	if err := binding.Validate(); err != nil {
		t.Fatalf("reproduced binding should validate: %v", err)
	}

	counterexample := ObserveExecutionEnvelopeProvenanceReverse(ExecutionEnvelopeReverseObservationInput{
		ObservedStatus:         "ready",
		ExpectedStatus:         "hold",
		ObservedEvidenceDigest: "observed",
		ExpectedEvidenceDigest: "expected",
		NonAuthorizing:         true,
	})
	counterexampleBinding := BindExecutionEnvelopeReverseObservation(counterexample)
	if counterexampleBinding.Status != "counterexample" || counterexampleBinding.FirstMismatch != "status" {
		t.Fatalf("counterexample binding lost mismatch: %#v", counterexampleBinding)
	}
	if err := counterexampleBinding.Validate(); err != nil {
		t.Fatalf("counterexample binding should validate: %v", err)
	}

	unknown := BindExecutionEnvelopeReverseObservation(ObserveExecutionEnvelopeProvenanceReverse(ExecutionEnvelopeReverseObservationInput{
		NonAuthorizing: true,
	}))
	if unknown.Status != "UNKNOWN" || unknown.FirstMismatch != "status" || unknown.ObservationDigest == "" {
		t.Fatalf("unknown binding lost first mismatch: %#v", unknown)
	}
	if err := unknown.Validate(); err != nil {
		t.Fatalf("unknown binding should validate: %v", err)
	}

	tampered := binding
	tampered.ObservationDigest = "tampered"
	if err := tampered.Validate(); err == nil {
		t.Fatal("tampered reverse observation must fail validation")
	}
}
