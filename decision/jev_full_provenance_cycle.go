package decision

import "strings"

// ExecutionEnvelopeJEVFullProvenanceCycleInput feeds a complete source chain
// and feedback aggregation into the existing improvement-cycle observation.
type ExecutionEnvelopeJEVFullProvenanceCycleInput struct {
	Provenance     ExecutionEnvelopeFullProvenanceSourceBinding
	ChangePlanDigest string
	Feedback       JEVImprovementFeedbackAggregation
	NonAuthorizing bool
}

// ObserveJEVImprovementCycleFromFullProvenance closes the source-to-cycle
// loop without executing, applying, or authorizing the proposed change.
func ObserveJEVImprovementCycleFromFullProvenance(input ExecutionEnvelopeJEVFullProvenanceCycleInput) JEVImprovementCycleObservation {
	output := JEVImprovementCycleObservation{
		Status: "UNKNOWN", NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.Provenance.NonAuthorizing || !input.Feedback.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		return output
	}
	if !input.Provenance.NonExecuting {
		output.MissingStage = "execution-boundary"
		return output
	}
	if input.Provenance.Status != "complete" {
		output.MissingStage = input.Provenance.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "full-provenance"
		}
		return output
	}
	stages := []struct {
		name  string
		value string
	}{
		{"declaration", input.Provenance.DeclarationDigest},
		{"ir", input.Provenance.IRDigest},
		{"generation", input.Provenance.GenerationDigest},
		{"reverse-observation", input.Provenance.ReverseObservationDigest},
		{"metric", input.Provenance.MetricDigest},
		{"change-plan", input.ChangePlanDigest},
	}
	for _, stage := range stages {
		if strings.TrimSpace(stage.value) == "" {
			output.MissingStage = stage.name
			return output
		}
	}
	return ObserveJEVImprovementFeedbackCycle(ExecutionEnvelopeJEVImprovementFeedbackCycleInput{
		DeclarationDigest:        input.Provenance.DeclarationDigest,
		IRDigest:                 input.Provenance.IRDigest,
		GenerationDigest:         input.Provenance.GenerationDigest,
		ReverseObservationDigest: input.Provenance.ReverseObservationDigest,
		MetricDigest:              input.Provenance.MetricDigest,
		ChangePlanDigest:         input.ChangePlanDigest,
		Feedback:                 input.Feedback,
		NonAuthorizing:           true,
	})
}
