package decision

import "testing"

func TestClassifyExecutionEnvelopeCounterexampleRequiresReview(t *testing.T) {
	result := ClassifyExecutionEnvelopeCounterexample(ExecutionEnvelopeCounterexampleDispositionInput{
		Status:         "counterexample",
		FirstMismatch:  "status",
		EvidenceDigest: "observed",
		NonAuthorizing: true,
	})
	if result.Status != "review-required" || !result.ReviewRequired || result.MismatchStage != "status" || !result.NonAuthorizing {
		t.Fatalf("counterexample did not require review: %#v", result)
	}
}

func TestClassifyExecutionEnvelopeCounterexampleKeepsReproductionObservationOnly(t *testing.T) {
	result := ClassifyExecutionEnvelopeCounterexample(ExecutionEnvelopeCounterexampleDispositionInput{
		Status:         "reproduced",
		EvidenceDigest: "observed",
		NonAuthorizing: true,
	})
	if result.Status != "observation-only" || result.ReviewRequired || result.EvidenceDigest != "observed" || !result.NonAuthorizing {
		t.Fatalf("reproduction became authorizing: %#v", result)
	}
}

func TestClassifyExecutionEnvelopeCounterexampleHoldsUnknown(t *testing.T) {
	result := ClassifyExecutionEnvelopeCounterexample(ExecutionEnvelopeCounterexampleDispositionInput{
		Status:         "UNKNOWN",
		NonAuthorizing: true,
	})
	if result.Status != "hold" || result.MismatchStage != "reverse-observation" || !result.NonAuthorizing {
		t.Fatalf("unknown observation escaped hold: %#v", result)
	}
}

func TestClassifyExecutionEnvelopeCounterexampleRejectsAuthorizationClaim(t *testing.T) {
	result := ClassifyExecutionEnvelopeCounterexample(ExecutionEnvelopeCounterexampleDispositionInput{
		Status:         "counterexample",
		FirstMismatch:  "status",
		EvidenceDigest: "observed",
		NonAuthorizing: false,
	})
	if result.Status != "UNKNOWN" || result.MismatchStage != "authorization-boundary" || result.NonAuthorizing {
		t.Fatalf("authorization claim escaped disposition: %#v", result)
	}
}
