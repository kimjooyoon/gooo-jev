package decision

// DecisionConfidenceEvidenceLedgerMetric is a one-hot observation of ledger integrity.
// It is telemetry only and cannot authorize execution or change application.
type DecisionConfidenceEvidenceLedgerMetric struct {
	EntryDigest    string `json:"entry_digest"`
	VerifiedCount  uint64 `json:"verified_count"`
	BrokenCount    uint64 `json:"broken_count"`
	Status         string `json:"status"`
	NonAuthorizing bool   `json:"non_authorizing"`
}

func ObserveDecisionConfidenceEvidenceLedgerMetric(entry DecisionConfidenceEvidenceLedgerEntry) DecisionConfidenceEvidenceLedgerMetric {
	if err := VerifyDecisionConfidenceEvidenceLedgerEntry(entry); err != nil {
		return DecisionConfidenceEvidenceLedgerMetric{
			EntryDigest:    entry.EntryDigest,
			BrokenCount:    1,
			Status:         "unknown",
			NonAuthorizing: true,
		}
	}
	return DecisionConfidenceEvidenceLedgerMetric{
		EntryDigest:    entry.EntryDigest,
		VerifiedCount:  1,
		Status:         "verified",
		NonAuthorizing: true,
	}
}
