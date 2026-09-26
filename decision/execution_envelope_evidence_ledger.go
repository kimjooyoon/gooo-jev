package decision

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// ExecutionEnvelopeEvidenceLedgerInput binds every provenance stage needed to
// seal one self-improvement observation.
type ExecutionEnvelopeEvidenceLedgerInput struct {
	DeclarationDigest       string `json:"declaration_digest"`
	IRDigest                string `json:"ir_digest"`
	GenerationDigest        string `json:"generation_digest"`
	ReverseObservationDigest string `json:"reverse_observation_digest"`
	MetricDigest            string `json:"metric_digest"`
	FeedbackDigest          string `json:"feedback_digest"`
	OutcomeDigest           string `json:"outcome_digest"`
	NonAuthorizing          bool   `json:"non_authorizing"`
}

// ExecutionEnvelopeEvidenceLedger records an ordered, content-addressed
// provenance boundary without claiming that the observation is an improvement.
type ExecutionEnvelopeEvidenceLedger struct {
	Status         string `json:"status"`
	LedgerDigest   string `json:"ledger_digest"`
	MissingStage   string `json:"missing_stage"`
	NonAuthorizing bool   `json:"non_authorizing"`
}

// SealExecutionEnvelopeEvidenceLedger fails closed at the first missing stage
// and hashes stage names with values to prevent ambiguous concatenation.
func SealExecutionEnvelopeEvidenceLedger(input ExecutionEnvelopeEvidenceLedgerInput) ExecutionEnvelopeEvidenceLedger {
	output := ExecutionEnvelopeEvidenceLedger{Status: "UNKNOWN", NonAuthorizing: true}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		output.MissingStage = "authorization-boundary"
		return output
	}
	stages := []struct {
		name  string
		value string
	}{
		{name: "declaration", value: input.DeclarationDigest},
		{name: "ir", value: input.IRDigest},
		{name: "generation", value: input.GenerationDigest},
		{name: "reverse_observation", value: input.ReverseObservationDigest},
		{name: "metric", value: input.MetricDigest},
		{name: "feedback", value: input.FeedbackDigest},
		{name: "outcome", value: input.OutcomeDigest},
	}
	parts := make([]string, 0, len(stages)*2)
	for _, stage := range stages {
		if stage.value == "" {
			output.MissingStage = stage.name
			return output
		}
		parts = append(parts, stage.name, stage.value)
	}
	digest := sha256.Sum256([]byte(strings.Join(parts, "\x00")))
	output.Status = "sealed"
	output.LedgerDigest = hex.EncodeToString(digest[:])
	return output
}
