package decision

import "testing"

func TestClassifyExecutionEnvelopeEvidenceLedgerCounterexampleRequiresReview(t *testing.T) {
	result := ClassifyExecutionEnvelopeEvidenceLedgerCounterexample(ExecutionEnvelopeEvidenceLedgerCounterexampleInput{
		Status:         "counterexample",
		FirstMismatch:  "metric",
		LedgerDigest:   "ledger",
		NonAuthorizing: true,
	})
	if result.Status != "review-required" || !result.ReviewRequired || result.MismatchStage != "metric" || result.LedgerDigest != "ledger" || !result.NonAuthorizing {
		t.Fatalf("ledger counterexample did not require review: %#v", result)
	}
}

func TestClassifyExecutionEnvelopeEvidenceLedgerCounterexampleKeepsReproductionObservationOnly(t *testing.T) {
	result := ClassifyExecutionEnvelopeEvidenceLedgerCounterexample(ExecutionEnvelopeEvidenceLedgerCounterexampleInput{
		Status:         "reproduced",
		LedgerDigest:   "ledger",
		NonAuthorizing: true,
	})
	if result.Status != "observation-only" || result.ReviewRequired || result.LedgerDigest != "ledger" || !result.NonAuthorizing {
		t.Fatalf("ledger reproduction became review-required: %#v", result)
	}
}

func TestClassifyExecutionEnvelopeEvidenceLedgerCounterexampleHoldsUnknown(t *testing.T) {
	result := ClassifyExecutionEnvelopeEvidenceLedgerCounterexample(ExecutionEnvelopeEvidenceLedgerCounterexampleInput{
		Status:         "UNKNOWN",
		NonAuthorizing: true,
	})
	if result.Status != "hold" || result.MismatchStage != "ledger-reverse-observation" || !result.NonAuthorizing {
		t.Fatalf("unknown ledger observation escaped hold: %#v", result)
	}
}

func TestClassifyExecutionEnvelopeEvidenceLedgerCounterexampleRejectsAuthorizationClaim(t *testing.T) {
	result := ClassifyExecutionEnvelopeEvidenceLedgerCounterexample(ExecutionEnvelopeEvidenceLedgerCounterexampleInput{
		Status:         "counterexample",
		FirstMismatch:  "metric",
		LedgerDigest:   "ledger",
		NonAuthorizing: false,
	})
	if result.Status != "UNKNOWN" || result.MismatchStage != "authorization-boundary" || result.NonAuthorizing {
		t.Fatalf("authorization claim escaped ledger counterexample: %#v", result)
	}
}
