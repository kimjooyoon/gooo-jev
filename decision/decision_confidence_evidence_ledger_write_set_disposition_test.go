package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerWriteSetDispositionIsFailClosed(t *testing.T) {
	verified := DecisionConfidenceEvidenceLedgerWriteSetIntegrity{
		ObservedWriteSetDigest: "write-set",
		ExpectedWriteSetDigest: "write-set",
		Status:                 "verified",
		NonAuthorizing:         true,
	}
	disposition, err := DeriveDecisionConfidenceEvidenceLedgerWriteSetDisposition(verified)
	if err != nil {
		t.Fatalf("derive verified disposition: %v", err)
	}
	if disposition.Status != "ready-for-external-apply-review" || !disposition.NonAuthorizing {
		t.Fatalf("unexpected verified disposition: %#v", disposition)
	}

	mismatch := verified
	mismatch.Status = "mismatch"
	disposition, err = DeriveDecisionConfidenceEvidenceLedgerWriteSetDisposition(mismatch)
	if err != nil {
		t.Fatalf("derive mismatch disposition: %v", err)
	}
	if disposition.Status != "abort" || !disposition.NonAuthorizing {
		t.Fatalf("mismatch was not aborted: %#v", disposition)
	}
}
