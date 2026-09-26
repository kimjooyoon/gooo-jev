package decision

// DecisionConfidenceEvidenceLedgerWriteSetObservation records an immutable
// write-set digest before any external change application.
type DecisionConfidenceEvidenceLedgerWriteSetObservation struct {
	ReviewSignalDigest string `json:"review_signal_digest"`
	WriteSetDigest     string `json:"write_set_digest"`
	Status             string `json:"status"`
	NonAuthorizing     bool   `json:"non_authorizing"`
}

func ObserveDecisionConfidenceEvidenceLedgerWriteSet(signal DecisionConfidenceEvidenceLedgerChangePlanReviewSignal, writeSetDigest string) DecisionConfidenceEvidenceLedgerWriteSetObservation {
	observation := DecisionConfidenceEvidenceLedgerWriteSetObservation{
		ReviewSignalDigest: signal.FeedbackDigest,
		WriteSetDigest:     writeSetDigest,
		Status:             "unknown",
		NonAuthorizing:     true,
	}
	if signal.Status == "ready-for-external-change-plan-review" && writeSetDigest != "" {
		observation.Status = "write-set-observed"
	}
	return observation
}
