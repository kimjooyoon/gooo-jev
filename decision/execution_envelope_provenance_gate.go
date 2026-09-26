package decision

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// ExecutionEnvelopeProvenanceGateInput binds the stages required before a
// self-improvement decision can be observed as ready.
type ExecutionEnvelopeProvenanceGateInput struct {
	DeclarationDigest       string `json:"declaration_digest"`
	IRDigest                string `json:"ir_digest"`
	GenerationDigest        string `json:"generation_digest"`
	ReverseObservationDigest string `json:"reverse_observation_digest"`
	MetricDigest            string `json:"metric_digest"`
	NonAuthorizing          bool   `json:"non_authorizing"`
}

// ExecutionEnvelopeProvenanceGateOutput records a fail-closed stage result.
type ExecutionEnvelopeProvenanceGateOutput struct {
	Status         string `json:"status"`
	MissingStage   string `json:"missing_stage"`
	EvidenceDigest string `json:"evidence_digest"`
	NonAuthorizing bool   `json:"non_authorizing"`
}

// EvaluateExecutionEnvelopeProvenanceGate never infers readiness from a
// partial chain and never turns provenance evidence into authorization.
func EvaluateExecutionEnvelopeProvenanceGate(input ExecutionEnvelopeProvenanceGateInput) ExecutionEnvelopeProvenanceGateOutput {
	output := ExecutionEnvelopeProvenanceGateOutput{Status: "UNKNOWN", NonAuthorizing: true}
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
	}
	parts := make([]string, 0, len(stages))
	for _, stage := range stages {
		if stage.value == "" {
			output.MissingStage = stage.name
			return output
		}
		parts = append(parts, stage.name, stage.value)
	}
	digest := sha256.Sum256([]byte(strings.Join(parts, "\x00")))
	output.Status = "ready"
	output.EvidenceDigest = hex.EncodeToString(digest[:])
	return output
}
