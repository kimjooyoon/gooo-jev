package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeReverseObservationInput compares an observed provenance
// boundary with the expected boundary without inferring missing evidence.
type ExecutionEnvelopeReverseObservationInput struct {
	ObservedStatus         string `json:"observed_status"`
	ExpectedStatus         string `json:"expected_status"`
	ObservedMissingStage   string `json:"observed_missing_stage"`
	ExpectedMissingStage   string `json:"expected_missing_stage"`
	ObservedEvidenceDigest string `json:"observed_evidence_digest"`
	ExpectedEvidenceDigest string `json:"expected_evidence_digest"`
	NonAuthorizing         bool   `json:"non_authorizing"`
}

// ExecutionEnvelopeReverseObservationOutput records the first exact
// divergence, or reproduced when every applicable boundary agrees.
type ExecutionEnvelopeReverseObservationOutput struct {
	Status         string `json:"status"`
	FirstMismatch  string `json:"first_mismatch"`
	EvidenceDigest string `json:"evidence_digest"`
	NonAuthorizing bool   `json:"non_authorizing"`
}

// ObserveExecutionEnvelopeProvenanceReverse preserves UNKNOWN for incomplete
// observations and never treats a reproduction as authorization.
func ObserveExecutionEnvelopeProvenanceReverse(input ExecutionEnvelopeReverseObservationInput) ExecutionEnvelopeReverseObservationOutput {
	output := ExecutionEnvelopeReverseObservationOutput{Status: "UNKNOWN", NonAuthorizing: true}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		output.FirstMismatch = "authorization-boundary"
		return output
	}
	checks := []struct {
		name     string
		observed string
		expected string
	}{
		{name: "status", observed: input.ObservedStatus, expected: input.ExpectedStatus},
		{name: "missing_stage", observed: input.ObservedMissingStage, expected: input.ExpectedMissingStage},
		{name: "evidence_digest", observed: input.ObservedEvidenceDigest, expected: input.ExpectedEvidenceDigest},
	}
	for _, check := range checks {
		if check.name == "missing_stage" && check.observed == "" && check.expected == "" && input.ObservedStatus == "ready" && input.ExpectedStatus == "ready" {
			continue
		}
		if check.observed == "" || check.expected == "" {
			output.FirstMismatch = check.name
			return output
		}
		if check.observed != check.expected {
			output.Status = "counterexample"
			output.FirstMismatch = check.name
			output.EvidenceDigest = input.ObservedEvidenceDigest
			return output
		}
	}
	output.Status = "reproduced"
	output.EvidenceDigest = input.ObservedEvidenceDigest
	return output
}

// ExecutionEnvelopeReverseObservationBinding seals the observation result
// itself so a reverse observation can be used as provenance evidence later.
type ExecutionEnvelopeReverseObservationBinding struct {
	Status                 string
	FirstMismatch          string
	ObservedEvidenceDigest string
	ObservationDigest      string
	NonAuthorizing         bool
}

// BindExecutionEnvelopeReverseObservation content-addresses a reverse result
// without turning reproduction or counterexamples into authorization.
func BindExecutionEnvelopeReverseObservation(output ExecutionEnvelopeReverseObservationOutput) ExecutionEnvelopeReverseObservationBinding {
	binding := ExecutionEnvelopeReverseObservationBinding{Status: "UNKNOWN", NonAuthorizing: true}
	if !output.NonAuthorizing {
		binding.NonAuthorizing = false
		binding.FirstMismatch = "authorization-boundary"
		return binding
	}
	switch output.Status {
	case "UNKNOWN":
		if strings.TrimSpace(output.FirstMismatch) == "" || strings.TrimSpace(output.EvidenceDigest) != "" {
			binding.FirstMismatch = "observation-evidence"
			return binding
		}
	case "reproduced", "counterexample":
		if strings.TrimSpace(output.EvidenceDigest) == "" || (output.Status == "counterexample" && strings.TrimSpace(output.FirstMismatch) == "") {
			binding.FirstMismatch = "observation-evidence"
			return binding
		}
	default:
		binding.FirstMismatch = "observation-status"
		return binding
	}
	digest, err := Digest(struct {
		Status         string
		FirstMismatch  string
		EvidenceDigest string
	}{
		Status:         output.Status,
		FirstMismatch:  output.FirstMismatch,
		EvidenceDigest: output.EvidenceDigest,
	})
	if err != nil {
		binding.FirstMismatch = "observation-digest"
		return binding
	}
	binding.Status = output.Status
	binding.FirstMismatch = output.FirstMismatch
	binding.ObservedEvidenceDigest = output.EvidenceDigest
	binding.ObservationDigest = digest
	return binding
}

// Validate replays the observation digest and keeps incomplete bindings
// distinguishable from reproduced or counterexample results.
func (binding ExecutionEnvelopeReverseObservationBinding) Validate() error {
	if !binding.NonAuthorizing || strings.TrimSpace(binding.ObservationDigest) == "" {
		return fmt.Errorf("reverse observation binding is incomplete")
	}
	switch binding.Status {
	case "UNKNOWN":
		if strings.TrimSpace(binding.FirstMismatch) == "" || strings.TrimSpace(binding.ObservedEvidenceDigest) != "" {
			return fmt.Errorf("unknown reverse observation binding has invalid evidence")
		}
	case "reproduced":
		if strings.TrimSpace(binding.FirstMismatch) != "" || strings.TrimSpace(binding.ObservedEvidenceDigest) == "" {
			return fmt.Errorf("reproduced reverse observation binding has invalid evidence")
		}
	case "counterexample":
		if strings.TrimSpace(binding.FirstMismatch) == "" || strings.TrimSpace(binding.ObservedEvidenceDigest) == "" {
			return fmt.Errorf("counterexample reverse observation binding has invalid evidence")
		}
	default:
		return fmt.Errorf("unsupported reverse observation status %q", binding.Status)
	}
	expected, err := Digest(struct {
		Status         string
		FirstMismatch  string
		EvidenceDigest string
	}{
		Status:         binding.Status,
		FirstMismatch:  binding.FirstMismatch,
		EvidenceDigest: binding.ObservedEvidenceDigest,
	})
	if err != nil {
		return err
	}
	if expected != binding.ObservationDigest {
		return fmt.Errorf("reverse observation binding digest mismatch")
	}
	return nil
}
