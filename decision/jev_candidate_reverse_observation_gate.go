package decision

import "strings"

// ExecutionEnvelopeJEVCandidateReverseObservationGateInput joins candidate
// admission evidence with a sealed reverse observation.
type ExecutionEnvelopeJEVCandidateReverseObservationGateInput struct {
	CandidateBinding ExecutionEnvelopeJEVSelfImprovementCandidateEvidenceGateBinding
	ReverseBinding   ExecutionEnvelopeReverseObservationBinding
	NonAuthorizing   bool
}

// ExecutionEnvelopeJEVCandidateReverseObservationGateBinding records the
// candidate and reverse evidence identities without calling either outcome a
// successful mutation.
type ExecutionEnvelopeJEVCandidateReverseObservationGateBinding struct {
	Status                   string
	CandidateDigest          string
	CandidateEvidenceDigest  string
	AdmissionDigest          string
	ReverseObservationStatus string
	ReverseObservationDigest string
	ReverseEvidenceDigest    string
	BindingDigest            string
	MissingStage             string
	NonExecuting             bool
	NonAuthorizing           bool
}

// BindExecutionEnvelopeJEVCandidateToReverseObservation requires a complete
// candidate evidence gate and a reproduced or counterexample reverse result.
func BindExecutionEnvelopeJEVCandidateToReverseObservation(input ExecutionEnvelopeJEVCandidateReverseObservationGateInput) ExecutionEnvelopeJEVCandidateReverseObservationGateBinding {
	output := ExecutionEnvelopeJEVCandidateReverseObservationGateBinding{
		Status: "UNKNOWN", NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.CandidateBinding.NonAuthorizing || !input.ReverseBinding.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		return output
	}
	if !input.CandidateBinding.NonExecuting {
		output.MissingStage = "execution-boundary"
		return output
	}
	if input.CandidateBinding.Status != "bound" || strings.TrimSpace(input.CandidateBinding.AdmissionDigest) == "" {
		output.MissingStage = input.CandidateBinding.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "candidate-evidence-gate"
		}
		return output
	}
	if err := input.ReverseBinding.Validate(); err != nil {
		output.MissingStage = "reverse-observation"
		return output
	}
	if input.ReverseBinding.Status == "UNKNOWN" {
		output.MissingStage = input.ReverseBinding.FirstMismatch
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "reverse-observation"
		}
		return output
	}
	if input.ReverseBinding.Status != "reproduced" && input.ReverseBinding.Status != "counterexample" {
		output.MissingStage = "reverse-observation-status"
		return output
	}
	bindingDigest, err := Digest(struct {
		CandidateDigest          string
		CandidateEvidenceDigest  string
		AdmissionDigest          string
		ReverseObservationStatus string
		ReverseObservationDigest string
		ReverseEvidenceDigest    string
	}{
		CandidateDigest:          input.CandidateBinding.CandidateDigest,
		CandidateEvidenceDigest:  input.CandidateBinding.CandidateEvidenceDigest,
		AdmissionDigest:          input.CandidateBinding.AdmissionDigest,
		ReverseObservationStatus: input.ReverseBinding.Status,
		ReverseObservationDigest: input.ReverseBinding.ObservationDigest,
		ReverseEvidenceDigest:    input.ReverseBinding.ObservedEvidenceDigest,
	})
	if err != nil {
		output.MissingStage = "candidate-reverse-binding-digest"
		return output
	}
	output.Status = "bound"
	output.CandidateDigest = input.CandidateBinding.CandidateDigest
	output.CandidateEvidenceDigest = input.CandidateBinding.CandidateEvidenceDigest
	output.AdmissionDigest = input.CandidateBinding.AdmissionDigest
	output.ReverseObservationStatus = input.ReverseBinding.Status
	output.ReverseObservationDigest = input.ReverseBinding.ObservationDigest
	output.ReverseEvidenceDigest = input.ReverseBinding.ObservedEvidenceDigest
	output.BindingDigest = bindingDigest
	return output
}
