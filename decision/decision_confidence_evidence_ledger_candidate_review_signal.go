package decision

// DecisionConfidenceEvidenceLedgerCandidateReviewSignal is an analysis-only
// bridge from external review evidence to a future candidate review boundary.
type DecisionConfidenceEvidenceLedgerCandidateReviewSignal struct {
	ReviewOutcomeDigest string `json:"review_outcome_digest"`
	Status              string `json:"status"`
	NonAuthorizing      bool   `json:"non_authorizing"`
}

func DeriveDecisionConfidenceEvidenceLedgerCandidateReviewSignal(outcome DecisionConfidenceEvidenceLedgerExternalReviewOutcome) (DecisionConfidenceEvidenceLedgerCandidateReviewSignal, error) {
	digest, err := Digest(outcome)
	if err != nil {
		return DecisionConfidenceEvidenceLedgerCandidateReviewSignal{}, err
	}
	signal := DecisionConfidenceEvidenceLedgerCandidateReviewSignal{
		ReviewOutcomeDigest: digest,
		Status:              "hold",
		NonAuthorizing:      true,
	}
	switch outcome.Status {
	case "accepted-for-analysis":
		signal.Status = "ready-for-candidate-review"
	case "rejected":
		signal.Status = "rejected"
	}
	return signal, nil
}
