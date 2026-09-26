package decision

import "strings"

// ExecutionEnvelopeGoooRevisionCandidateInput connects full provenance,
// feedback direction, and revision source evidence.
type ExecutionEnvelopeGoooRevisionCandidateInput struct {
	Provenance           ExecutionEnvelopeFullProvenanceSourceBinding
	ChangePlanDigest     string
	Feedback             JEVImprovementFeedbackAggregation
	Directive            JEVImprovementDirectionDirective
	RevisionSource       string
	RevisionChangeDigest string
	NonAuthorizing       bool
}

// GenerateExecutionEnvelopeJEVRevisionCandidateFromGoooFullProvenance
// closes the source-to-candidate path without applying or authorizing a
// revision.
func GenerateExecutionEnvelopeJEVRevisionCandidateFromGoooFullProvenance(input ExecutionEnvelopeGoooRevisionCandidateInput) ExecutionEnvelopeJEVCycleRevisionCandidateGenerationBinding {
	output := ExecutionEnvelopeJEVCycleRevisionCandidateGenerationBinding{
		Status: "UNKNOWN", NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.Provenance.NonAuthorizing || !input.Feedback.NonAuthorizing || !input.Directive.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		return output
	}
	cycle := ObserveJEVImprovementCycleFromFullProvenance(ExecutionEnvelopeJEVFullProvenanceCycleInput{
		Provenance:      input.Provenance,
		ChangePlanDigest: input.ChangePlanDigest,
		Feedback:        input.Feedback,
		NonAuthorizing:  true,
	})
	if cycle.Status == "UNKNOWN" {
		output.MissingStage = cycle.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "improvement-cycle-observation"
		}
		return output
	}
	if err := input.Directive.Validate(); err != nil {
		output.MissingStage = "direction-directive"
		return output
	}
	if input.Directive.InputEvidenceDigest != input.Feedback.EvidenceDigest {
		output.MissingStage = "feedback-direction-binding"
		return output
	}
	return GenerateExecutionEnvelopeJEVRevisionCandidateFromCycle(ExecutionEnvelopeJEVCycleRevisionCandidateGenerationInput{
		Cycle:                cycle,
		Directive:            input.Directive,
		RevisionSource:       input.RevisionSource,
		RevisionChangeDigest: input.RevisionChangeDigest,
		NonAuthorizing:      true,
	})
}
