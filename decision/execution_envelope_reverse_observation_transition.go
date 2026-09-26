package decision

// ExecutionEnvelopeReverseObservationTransition records observed changes
// between two validated reverse-observation receipts.
type ExecutionEnvelopeReverseObservationTransition struct {
	Status             string
	Changed            bool
	Direction          string
	FirstMismatch      string
	ObservationDigest  string
	NonAuthorizing     bool
}

// CompareExecutionEnvelopeReverseObservationTransition compares structured
// reverse evidence without inferring improvement from a transition.
func CompareExecutionEnvelopeReverseObservationTransition(
	previous ReverseObservation,
	current ReverseObservation,
) ExecutionEnvelopeReverseObservationTransition {
	output := ExecutionEnvelopeReverseObservationTransition{
		Status: "UNKNOWN", NonAuthorizing: true,
	}
	if err := previous.Validate(); err != nil {
		output.FirstMismatch = "previous-reverse-observation"
		return output
	}
	if err := current.Validate(); err != nil {
		output.FirstMismatch = "current-reverse-observation"
		return output
	}
	output.ObservationDigest = current.ObservationDigest
	if previous.Status != current.Status {
		output.Status = "changed"
		output.Changed = true
		output.Direction = "status-transition"
		output.FirstMismatch = "reverse-observation-status"
		return output
	}
	if previous.MissingStage != current.MissingStage {
		output.Status = "changed"
		output.Changed = true
		output.Direction = "missing-stage-transition"
		output.FirstMismatch = "reverse-observation-missing-stage"
		return output
	}
	if previous.ObservedOutputDigest != current.ObservedOutputDigest {
		output.Status = "changed"
		output.Changed = true
		output.Direction = "output-transition"
		output.FirstMismatch = "reverse-observation-output"
		return output
	}
	if previous.VerifierDigest != current.VerifierDigest {
		output.Status = "changed"
		output.Changed = true
		output.Direction = "verifier-transition"
		output.FirstMismatch = "reverse-observation-verifier"
		return output
	}
	if previous.ObservationDigest != current.ObservationDigest {
		output.Status = "changed"
		output.Changed = true
		output.Direction = "digest-transition"
		output.FirstMismatch = "reverse-observation-digest"
		return output
	}
	output.Status = "unchanged"
	output.Direction = "stable"
	return output
}