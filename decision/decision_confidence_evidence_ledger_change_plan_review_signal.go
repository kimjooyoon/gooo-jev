package decision

// DecisionConfidenceEvidenceLedgerChangePlanReviewSignal is a non-executing
// boundary before any external change-plan review or authorization.
type DecisionConfidenceEvidenceLedgerChangePlanReviewSignal struct {
	FeedbackDigest string `json:"feedback_digest"`
	Status         string `json:"status"`
	NonAuthorizing bool   `json:"non_authorizing"`
}

func DeriveDecisionConfidenceEvidenceLedgerChangePlanReviewSignal(feedback DecisionConfidenceEvidenceLedgerFeedbackSignal) (DecisionConfidenceEvidenceLedgerChangePlanReviewSignal, error) {
	digest, err := Digest(feedback)
	if err != nil {
		return DecisionConfidenceEvidenceLedgerChangePlanReviewSignal{}, err
	}
	signal := DecisionConfidenceEvidenceLedgerChangePlanReviewSignal{
		FeedbackDigest: digest,
		Status:         "hold",
		NonAuthorizing: true,
	}
	switch feedback.Status {
	case "review-required":
		signal.Status = "ready-for-external-change-plan-review"
	case "observation-only":
		signal.Status = "observation-only"
	}
	return signal, nil
}
