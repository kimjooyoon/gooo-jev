package decision

// DecisionConfidenceEvidenceLedgerCandidateReviewReceipt is a content-addressed
// record of review evidence, not a candidate or authorization.
type DecisionConfidenceEvidenceLedgerCandidateReviewReceipt struct {
	SignalDigest        string `json:"signal_digest"`
	ReviewOutcomeDigest  string `json:"review_outcome_digest"`
	ReceiptDigest       string `json:"receipt_digest"`
	Status              string `json:"status"`
	NonAuthorizing      bool   `json:"non_authorizing"`
}

type decisionConfidenceEvidenceLedgerCandidateReviewReceiptIdentity struct {
	SignalDigest       string `json:"signal_digest"`
	ReviewOutcomeDigest string `json:"review_outcome_digest"`
}

func BuildDecisionConfidenceEvidenceLedgerCandidateReviewReceipt(signal DecisionConfidenceEvidenceLedgerCandidateReviewSignal, outcome DecisionConfidenceEvidenceLedgerExternalReviewOutcome) (DecisionConfidenceEvidenceLedgerCandidateReviewReceipt, error) {
	identity := decisionConfidenceEvidenceLedgerCandidateReviewReceiptIdentity{
		SignalDigest:        signal.ReviewOutcomeDigest,
		ReviewOutcomeDigest: outcome.HandoffDigest,
	}
	digest, err := Digest(identity)
	if err != nil {
		return DecisionConfidenceEvidenceLedgerCandidateReviewReceipt{}, err
	}
	receipt := DecisionConfidenceEvidenceLedgerCandidateReviewReceipt{
		SignalDigest:       signal.ReviewOutcomeDigest,
		ReviewOutcomeDigest: outcome.HandoffDigest,
		ReceiptDigest:      digest,
		Status:             "hold",
		NonAuthorizing:     true,
	}
	if signal.Status == "ready-for-candidate-review" && outcome.Status == "accepted-for-analysis" {
		receipt.Status = "ready-for-candidate-review"
	} else if outcome.Status == "rejected" {
		receipt.Status = "rejected"
	}
	return receipt, nil
}
