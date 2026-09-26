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
	if result.Status != "reproduced" || result.FirstMismatch != "missing_stage" || !result.NonAuthorizing {
		t.Fatalf("incomplete reverse boundary was not preserved: %#v", result)
	}
}

func TestObserveExecutionEnvelopeProvenanceReverseReportsFirstCounterexample(t *testing.T) {
	result := ObserveExecutionEnvelopeProvenanceReverse(ExecutionEnvelopeReverseObservationInput{
		ObservedStatus:         "ready",
		ExpectedStatus:         "hold",
		ObservedMissingStage:   "",
		ExpectedMissingStage:   "",
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
