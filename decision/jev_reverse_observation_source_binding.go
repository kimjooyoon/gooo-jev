package decision

import "strings"

// ExecutionEnvelopeReverseObservationSourceInput carries the exact source
// representation and typed fields used for a non-authorizing reverse check.
type ExecutionEnvelopeReverseObservationSourceInput struct {
	ObservedStatus         string
	ExpectedStatus         string
	ObservedMissingStage   string
	ExpectedMissingStage   string
	ObservedEvidenceDigest string
	ExpectedEvidenceDigest string
	SourceText             string
	NonAuthorizing         bool
}

// ExecutionEnvelopeReverseObservationSourceBinding seals both the source
// origin and the typed reverse-observation result.
type ExecutionEnvelopeReverseObservationSourceBinding struct {
	Status            string
	FirstMismatch     string
	SourceDigest      string
	ObservationDigest string
	NonAuthorizing    bool
}

// BindExecutionEnvelopeReverseObservationFromSource refuses to infer a
// reproduced observation from an absent source or an unauthorized boundary.
func BindExecutionEnvelopeReverseObservationFromSource(input ExecutionEnvelopeReverseObservationSourceInput) ExecutionEnvelopeReverseObservationSourceBinding {
	output := ExecutionEnvelopeReverseObservationSourceBinding{
		Status: "UNKNOWN", NonAuthorizing: true,
	}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		output.FirstMismatch = "authorization-boundary"
		return output
	}
	if strings.TrimSpace(input.SourceText) == "" {
		output.FirstMismatch = "reverse-observation-source"
		return output
	}
	sourceDigest, err := Digest(struct {
		ObservedStatus         string
		ExpectedStatus         string
		ObservedMissingStage   string
		ExpectedMissingStage   string
		ObservedEvidenceDigest string
		ExpectedEvidenceDigest string
		SourceText             string
	}{
		ObservedStatus:         input.ObservedStatus,
		ExpectedStatus:         input.ExpectedStatus,
		ObservedMissingStage:   input.ObservedMissingStage,
		ExpectedMissingStage:   input.ExpectedMissingStage,
		ObservedEvidenceDigest: input.ObservedEvidenceDigest,
		ExpectedEvidenceDigest: input.ExpectedEvidenceDigest,
		SourceText:             input.SourceText,
	})
	if err != nil {
		output.FirstMismatch = "reverse-observation-source-digest"
		return output
	}
	observation := ObserveExecutionEnvelopeProvenanceReverse(ExecutionEnvelopeReverseObservationInput{
		ObservedStatus:         input.ObservedStatus,
		ExpectedStatus:         input.ExpectedStatus,
		ObservedMissingStage:   input.ObservedMissingStage,
		ExpectedMissingStage:   input.ExpectedMissingStage,
		ObservedEvidenceDigest: input.ObservedEvidenceDigest,
		ExpectedEvidenceDigest: input.ExpectedEvidenceDigest,
		NonAuthorizing:         true,
	})
	binding := BindExecutionEnvelopeReverseObservation(observation)
	output.Status = binding.Status
	output.FirstMismatch = binding.FirstMismatch
	output.SourceDigest = sourceDigest
	output.ObservationDigest = binding.ObservationDigest
	return output
}
