package decision

import (
	"crypto/sha256"
	"encoding/hex"
)

// ExecutionEnvelopeReviewOutcomeInput binds a review handoff to an explicit
// outcome evidence record.
type ExecutionEnvelopeReviewOutcomeInput struct {
	HandoffStatus        string `json:"handoff_status"`
	HandoffDigest        string `json:"handoff_digest"`
	ReviewRequired       bool   `json:"review_required"`
	OutcomeStatus        string `json:"outcome_status"`
	OutcomeEvidenceDigest string `json:"outcome_evidence_digest"`
	NonAuthorizing       bool   `json:"non_authorizing"`
}

// ExecutionEnvelopeReviewOutcome records a bounded external review result.
type ExecutionEnvelopeReviewOutcome struct {
	Status         string `json:"status"`
	OutcomeDigest  string `json:"outcome_digest"`
	MissingStage   string `json:"missing_stage"`
	NonAuthorizing bool   `json:"non_authorizing"`
}

// RecordExecutionEnvelopeReviewOutcome accepts only explicit outcomes for a
// valid handoff and never treats the result as authorization.
func RecordExecutionEnvelopeReviewOutcome(input ExecutionEnvelopeReviewOutcomeInput) ExecutionEnvelopeReviewOutcome {
	output := ExecutionEnvelopeReviewOutcome{Status: "UNKNOWN", NonAuthorizing: true}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		output.MissingStage = "authorization-boundary"
		return output
	}
	if input.HandoffStatus != "handoff-required" || !input.ReviewRequired || input.HandoffDigest == "" {
		output.Status = "hold"
		output.MissingStage = "review-handoff"
		return output
	}
	if input.OutcomeStatus == "" || input.OutcomeEvidenceDigest == "" {
		output.MissingStage = "review-outcome"
		return output
	}
	switch input.OutcomeStatus {
	case "accepted", "rejected", "hold":
		digest := sha256.Sum256([]byte(input.HandoffDigest + "\x00" + input.OutcomeStatus + "\x00" + input.OutcomeEvidenceDigest))
		output.Status = input.OutcomeStatus
		output.OutcomeDigest = hex.EncodeToString(digest[:])
		return output
	default:
		output.MissingStage = "review-outcome"
		return output
	}
}
