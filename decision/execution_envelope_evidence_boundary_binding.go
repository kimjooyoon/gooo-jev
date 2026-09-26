package decision

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// ExecutionEnvelopeEvidenceBoundaryBindingInput joins provenance and
// capability-boundary evidence without authorizing execution.
type ExecutionEnvelopeEvidenceBoundaryBindingInput struct {
	LedgerDigest    string
	BoundaryDigest  string
	NonAuthorizing  bool
}

// ExecutionEnvelopeEvidenceBoundaryBinding records a content-addressed link
// between the evidence ledger and its execution boundary.
type ExecutionEnvelopeEvidenceBoundaryBinding struct {
	Status          string
	BindingDigest   string
	MissingStage    string
	NonAuthorizing  bool
}

// BindExecutionEnvelopeEvidenceBoundary fails closed when either side of the
// provenance boundary is missing.
func BindExecutionEnvelopeEvidenceBoundary(input ExecutionEnvelopeEvidenceBoundaryBindingInput) ExecutionEnvelopeEvidenceBoundaryBinding {
	output := ExecutionEnvelopeEvidenceBoundaryBinding{Status: "UNKNOWN", NonAuthorizing: true}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		output.MissingStage = "authorization-boundary"
		return output
	}
	if input.LedgerDigest == "" {
		output.MissingStage = "evidence-ledger"
		return output
	}
	if input.BoundaryDigest == "" {
		output.MissingStage = "capability-boundary"
		return output
	}
	parts := []string{"evidence-ledger", input.LedgerDigest, "capability-boundary", input.BoundaryDigest}
	digest := sha256.Sum256([]byte(strings.Join(parts, string(rune(0)))))
	output.Status = "bound"
	output.BindingDigest = hex.EncodeToString(digest[:])
	return output
}