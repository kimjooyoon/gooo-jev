package decision

import "testing"

func TestSealExecutionEnvelopeEvidenceLedgerPreservesFirstMissingStage(t *testing.T) {
	result := SealExecutionEnvelopeEvidenceLedger(ExecutionEnvelopeEvidenceLedgerInput{
		IRDigest:       "ir",
		NonAuthorizing: true,
	})
	if result.Status != "UNKNOWN" || result.MissingStage != "declaration" || result.LedgerDigest != "" || !result.NonAuthorizing {
		t.Fatalf("partial ledger escaped fail-closed state: %#v", result)
	}
}

func TestSealExecutionEnvelopeEvidenceLedgerSealsOrderedStages(t *testing.T) {
	result := SealExecutionEnvelopeEvidenceLedger(ExecutionEnvelopeEvidenceLedgerInput{
		DeclarationDigest:        "declaration",
		IRDigest:                 "ir",
		GenerationDigest:         "generation",
		ReverseObservationDigest: "reverse",
		MetricDigest:             "metric",
		FeedbackDigest:           "feedback",
		OutcomeDigest:            "outcome",
		NonAuthorizing:           true,
	})
	if result.Status != "sealed" || result.MissingStage != "" || result.LedgerDigest == "" || !result.NonAuthorizing {
		t.Fatalf("complete ledger was not sealed: %#v", result)
	}
}

func TestSealExecutionEnvelopeEvidenceLedgerRejectsAuthorizationClaim(t *testing.T) {
	result := SealExecutionEnvelopeEvidenceLedger(ExecutionEnvelopeEvidenceLedgerInput{
		DeclarationDigest: "declaration",
		NonAuthorizing:    false,
	})
	if result.Status != "UNKNOWN" || result.MissingStage != "authorization-boundary" || result.LedgerDigest != "" || result.NonAuthorizing {
		t.Fatalf("authorization claim escaped ledger: %#v", result)
	}
}
