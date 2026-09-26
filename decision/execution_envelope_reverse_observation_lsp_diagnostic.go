package decision

// ProjectExecutionEnvelopeReverseObservationLSPDiagnostic bridges the
// validated reverse-observation model into the editor-facing projection.
func ProjectExecutionEnvelopeReverseObservationLSPDiagnostic(
	observation ReverseObservation,
	stageIndex int,
	evidencePrefixDigest string,
) ExecutionEnvelopeLSPDiagnostic {
	output := ExecutionEnvelopeLSPDiagnostic{
		Status: "UNKNOWN", MissingStageIndex: -1, NonAuthorizing: true,
	}
	if err := observation.Validate(); err != nil {
		output.Code = "reverse-observation-integrity"
		return output
	}
	if evidencePrefixDigest == "" {
		output.Code = "reverse-observation-prefix"
		return output
	}
	switch observation.Status {
	case ProvenanceObserved:
		if observation.MissingStage != "" || stageIndex != -1 {
			output.Code = "reverse-observation-shape"
			return output
		}
		output.Status = "clear"
		output.Severity = "info"
		output.Code = "provenance-complete"
		output.EvidencePrefixDigest = evidencePrefixDigest
	case ProvenanceUnknown:
		if observation.MissingStage == "" || stageIndex < 0 {
			output.Code = "reverse-observation-location"
			return output
		}
		output.Status = "publishable"
		output.Publishable = true
		output.Severity = "error"
		output.Code = "provenance-incomplete"
		output.MissingStageIndex = stageIndex
		output.EvidencePrefixDigest = evidencePrefixDigest
	case ProvenanceReview:
		if observation.MissingStage == "" || stageIndex < 0 {
			output.Code = "reverse-observation-location"
			return output
		}
		output.Status = "publishable"
		output.Publishable = true
		output.Severity = "warning"
		output.Code = "review-required"
		output.MissingStageIndex = stageIndex
		output.EvidencePrefixDigest = evidencePrefixDigest
	default:
		output.Code = "reverse-observation-status"
	}
	return output
}