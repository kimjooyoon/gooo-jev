package decision

// ImprovementReviewEvidenceGateInput binds candidate admission to its review
// receipt without allowing a review to grant execution.
type ImprovementReviewEvidenceGateInput struct {
	CandidateAdmissionStatus string
	CandidateDigest          string
	ReviewCandidateDigest    string
	ReviewerReference        string
	ReviewEvidenceDigest     string
	ReviewDecision           ImprovementReviewDecision
	ExecutionGranted         bool
	ReviewDigest             string
	NonAuthorizing           bool
}

// ImprovementReviewEvidenceGate classifies review disposition while keeping
// observation separate from execution.
type ImprovementReviewEvidenceGate struct {
	Status           string
	ObservationOnly  bool
	CandidateDigest  string
	ReviewDigest     string
	MissingStage     string
	ExecutionGranted bool
	NonAuthorizing   bool
}

// GateImprovementReviewEvidence requires a proposed candidate and a valid,
// digest-bound review receipt.
func GateImprovementReviewEvidence(input ImprovementReviewEvidenceGateInput) ImprovementReviewEvidenceGate {
	output := ImprovementReviewEvidenceGate{
		Status: "UNKNOWN", ExecutionGranted: false, NonAuthorizing: true,
	}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		output.MissingStage = "authorization-boundary"
		return output
	}
	if input.CandidateAdmissionStatus == "review" {
		output.Status = "review"
		output.MissingStage = "candidate-evidence-unknown"
		return output
	}
	if input.CandidateAdmissionStatus != "proposed" {
		output.MissingStage = "candidate-admission"
		return output
	}
	if input.CandidateDigest == "" {
		output.MissingStage = "candidate-digest"
		return output
	}
	if input.ReviewCandidateDigest != input.CandidateDigest {
		output.MissingStage = "candidate-review-binding"
		return output
	}
	if input.ExecutionGranted {
		output.NonAuthorizing = false
		output.MissingStage = "authorization-boundary"
		return output
	}
	review := ImprovementReviewReceipt{
		CandidateDigest:      input.ReviewCandidateDigest,
		ReviewerReference:    input.ReviewerReference,
		ReviewEvidenceDigest: input.ReviewEvidenceDigest,
		Decision:             input.ReviewDecision,
		ExecutionGranted:     input.ExecutionGranted,
		ReviewDigest:         input.ReviewDigest,
	}
	if err := review.Validate(); err != nil {
		output.MissingStage = "review-receipt"
		return output
	}
	output.CandidateDigest = input.CandidateDigest
	output.ReviewDigest = review.ReviewDigest
	switch review.Decision {
	case ReviewPassed:
		output.Status = "accepted-for-observation"
		output.ObservationOnly = true
	case ReviewFailed:
		output.Status = "rejected"
	case ReviewRequired:
		output.Status = "review"
		output.MissingStage = "review-required"
	case ReviewUnknown:
		output.MissingStage = "review-unknown"
	default:
		output.MissingStage = "review-decision"
	}
	return output
}