package decision

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// ExecutionEnvelopeReviewFeedbackInput binds a counterexample disposition to
// an explicit external review outcome without treating review as authorization.
type ExecutionEnvelopeReviewFeedbackInput struct {
	DispositionStatus       string `json:"disposition_status"`
	ReviewRequired          bool   `json:"review_required"`
	MismatchStage           string `json:"mismatch_stage"`
	EvidenceDigest          string `json:"evidence_digest"`
	ReviewStatus            string `json:"review_status"`
	ReviewerEvidenceDigest  string `json:"reviewer_evidence_digest"`
	NonAuthorizing          bool   `json:"non_authorizing"`
}

// ExecutionEnvelopeReviewFeedback records an explicit review result or a
// fail-closed hold when the result is incomplete.
type ExecutionEnvelopeReviewFeedback struct {
	Status         string `json:"status"`
	FeedbackDigest string `json:"feedback_digest"`
	ReviewRequired bool   `json:"review_required"`
	MissingStage   string `json:"missing_stage"`
	NonAuthorizing bool   `json:"non_authorizing"`
}

// RecordExecutionEnvelopeReviewFeedback preserves UNKNOWN for missing or
// unsupported review evidence and keeps every output non-authorizing.
func RecordExecutionEnvelopeReviewFeedback(input ExecutionEnvelopeReviewFeedbackInput) ExecutionEnvelopeReviewFeedback {
	output := ExecutionEnvelopeReviewFeedback{Status: "UNKNOWN", NonAuthorizing: true}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		output.MissingStage = "authorization-boundary"
		return output
	}
	if input.DispositionStatus == "observation-only" && !input.ReviewRequired {
		if input.EvidenceDigest == "" {
			output.MissingStage = "observation-evidence"
			return output
		}
		output.Status = "observation-only"
		output.FeedbackDigest = input.EvidenceDigest
		return output
	}
	if input.DispositionStatus != "review-required" || !input.ReviewRequired {
		output.Status = "hold"
		output.MissingStage = "counterexample-disposition"
		return output
	}
	if input.MismatchStage == "" || input.EvidenceDigest == "" || input.ReviewerEvidenceDigest == "" || input.ReviewStatus == "" {
		output.MissingStage = "review-outcome"
		return output
	}
	switch input.ReviewStatus {
	case "accepted", "rejected", "hold":
		output.Status = "reviewed"
		if input.ReviewStatus == "rejected" {
			output.Status = "review-rejected"
		}
		if input.ReviewStatus == "hold" {
			output.Status = "hold"
		}
		digest := sha256.Sum256([]byte(strings.Join([]string{input.EvidenceDigest, input.MismatchStage, input.ReviewStatus, input.ReviewerEvidenceDigest}, "\x00")))
		output.FeedbackDigest = hex.EncodeToString(digest[:])
		return output
	default:
		output.MissingStage = "review-outcome"
		return output
	}
}
