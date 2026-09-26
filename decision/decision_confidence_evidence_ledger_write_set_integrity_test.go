package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerWriteSetIntegrityIsFailClosed(t *testing.T) {
	observation := DecisionConfidenceEvidenceLedgerWriteSetObservation{
		ReviewSignalDigest: "review-digest",
		WriteSetDigest:     "write-set-digest",
		Status:             "write-set-observed",
		NonAuthorizing:     true,
	}
	verified := VerifyDecisionConfidenceEvidenceLedgerWriteSetIntegrity(observation, "write-set-digest")
	if verified.Status != "verified" || !verified.NonAuthorizing {
		t.Fatalf("unexpected verified write-set: %#v", verified)
	}
	mismatch := VerifyDecisionConfidenceEvidenceLedgerWriteSetIntegrity(observation, "other-write-set")
	if mismatch.Status != "mismatch" || !mismatch.NonAuthorizing {
		t.Fatalf("write-set mismatch was not preserved: %#v", mismatch)
	}
	unknown := VerifyDecisionConfidenceEvidenceLedgerWriteSetIntegrity(observation, "")
	if unknown.Status != "unknown" || !unknown.NonAuthorizing {
		t.Fatalf("missing expectation escaped unknown: %#v", unknown)
	}
}
