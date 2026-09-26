package decision

// DecisionConfidenceEvidenceLedgerWriteSetDisposition gates any later external
// apply review; it never applies a write-set itself.
type DecisionConfidenceEvidenceLedgerWriteSetDisposition struct {
	IntegrityDigest string `json:"integrity_digest"`
	Status         string `json:"status"`
	NonAuthorizing bool   `json:"non_authorizing"`
}

func DeriveDecisionConfidenceEvidenceLedgerWriteSetDisposition(integrity DecisionConfidenceEvidenceLedgerWriteSetIntegrity) (DecisionConfidenceEvidenceLedgerWriteSetDisposition, error) {
	digest, err := Digest(integrity)
	if err != nil {
		return DecisionConfidenceEvidenceLedgerWriteSetDisposition{}, err
	}
	disposition := DecisionConfidenceEvidenceLedgerWriteSetDisposition{
		IntegrityDigest: digest,
		Status:         "hold",
		NonAuthorizing: true,
	}
	switch integrity.Status {
	case "verified":
		disposition.Status = "ready-for-external-apply-review"
	case "mismatch":
		disposition.Status = "abort"
	}
	return disposition, nil
}
