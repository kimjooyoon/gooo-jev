package decision

import "testing"

func TestObserveExecutionEnvelopeEvidenceLedgerReverseReproducesSealedLedger(t *testing.T) {
	result := ObserveExecutionEnvelopeEvidenceLedgerReverse(ExecutionEnvelopeEvidenceLedgerReverseInput{
		ObservedStatus:       "sealed",
		ExpectedStatus:       "sealed",
		ObservedLedgerDigest: "ledger",
		ExpectedLedgerDigest: "ledger",
		NonAuthorizing:       true,
	})
	if result.Status != "reproduced" || result.FirstMismatch != "" || result.LedgerDigest != "ledger" || !result.NonAuthorizing {
		t.Fatalf("sealed ledger reproduction was not preserved: %#v", result)
	}
}

func TestObserveExecutionEnvelopeEvidenceLedgerReverseReportsFirstMismatch(t *testing.T) {
	result := ObserveExecutionEnvelopeEvidenceLedgerReverse(ExecutionEnvelopeEvidenceLedgerReverseInput{
		ObservedStatus:       "sealed",
		ExpectedStatus:       "UNKNOWN",
		ObservedLedgerDigest: "observed",
		ExpectedLedgerDigest: "expected",
		NonAuthorizing:       true,
	})
	if result.Status != "counterexample" || result.FirstMismatch != "status" || result.LedgerDigest != "observed" || !result.NonAuthorizing {
		t.Fatalf("ledger counterexample was not preserved: %#v", result)
	}
}

func TestObserveExecutionEnvelopeEvidenceLedgerReverseRejectsAuthorizationClaim(t *testing.T) {
	result := ObserveExecutionEnvelopeEvidenceLedgerReverse(ExecutionEnvelopeEvidenceLedgerReverseInput{
		ObservedStatus:       "sealed",
		ExpectedStatus:       "sealed",
		ObservedLedgerDigest: "ledger",
		ExpectedLedgerDigest: "ledger",
		NonAuthorizing:       false,
	})
	if result.Status != "UNKNOWN" || result.FirstMismatch != "authorization-boundary" || result.NonAuthorizing {
		t.Fatalf("authorization claim escaped ledger reverse audit: %#v", result)
	}
}
