package decision

// DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewHandoff is a
// non-executing handoff for review of receipt integrity counterexamples.
type DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewHandoff struct {
	CounterexampleDigest string `json:"counterexample_digest"`
	Status              string `json:"status"`
	NonAuthorizing      bool   `json:"non_authorizing"`
}

func DeriveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewHandoff(counterexample DecisionConfidenceEvidenceLedgerReceiptIntegrityCounterexample) (DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewHandoff, error) {
	digest, err := Digest(counterexample)
	if err != nil {
		return DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewHandoff{}, err
	}
	handoff := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewHandoff{
		CounterexampleDigest: digest,
		Status:              "hold",
		NonAuthorizing:      true,
	}
	switch counterexample.Status {
	case "counterexample":
		handoff.Status = "ready-for-external-review"
	case "no-counterexample":
		handoff.Status = "observation-only"
	}
	return handoff, nil
}
