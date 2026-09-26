package decision

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
