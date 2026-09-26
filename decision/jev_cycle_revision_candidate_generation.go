package decision

import "strings"

// ExecutionEnvelopeJEVCycleRevisionCandidateGenerationInput connects a
// validated improvement cycle to a revision candidate proposal.
type ExecutionEnvelopeJEVCycleRevisionCandidateGenerationInput struct {
	Cycle                JEVImprovementCycleObservation
	Directive            JEVImprovementDirectionDirective
	RevisionSource       string
	RevisionChangeDigest string
	NonAuthorizing       bool
}

// ExecutionEnvelopeJEVCycleRevisionCandidateGenerationBinding records a
// candidate whose change identity includes the cycle evidence.
type ExecutionEnvelopeJEVCycleRevisionCandidateGenerationBinding struct {
	Status                  string
	CycleStatus             string
	CycleEvidenceDigest     string
	CandidateStatus         string
	CandidateDigest         string
	CandidateEvidenceDigest string
	BoundRevisionChangeDigest string
	BindingDigest           string
	MissingStage            string
	NonExecuting            bool
	NonAuthorizing          bool
}

// GenerateExecutionEnvelopeJEVRevisionCandidateFromCycle creates only a
// non-executing, non-authorizing candidate evidence binding.
func GenerateExecutionEnvelopeJEVRevisionCandidateFromCycle(input ExecutionEnvelopeJEVCycleRevisionCandidateGenerationInput) ExecutionEnvelopeJEVCycleRevisionCandidateGenerationBinding {
	output := ExecutionEnvelopeJEVCycleRevisionCandidateGenerationBinding{
		Status: "UNKNOWN", NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.Cycle.NonAuthorizing || !input.Directive.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		return output
	}
	if !input.Cycle.NonExecuting || !input.Directive.NonExecuting {
		output.MissingStage = "execution-boundary"
		return output
	}
	if err := input.Cycle.Validate(); err != nil {
		output.MissingStage = input.Cycle.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "improvement-cycle-observation"
		}
		return output
	}
	switch input.Cycle.Status {
	case "stable-for-review", "needs-revision":
		// Both statuses produce evidence for review; neither authorizes a change.
	case "hold":
		output.MissingStage = "feedback-hold"
		return output
	default:
		output.MissingStage = "improvement-cycle-status"
		return output
	}
	if strings.TrimSpace(input.RevisionSource) == "" {
		output.MissingStage = "revision-source"
		return output
	}
	if strings.TrimSpace(input.RevisionChangeDigest) == "" {
		output.MissingStage = "revision-change"
		return output
	}
	boundChangeDigest, err := Digest(struct {
		RevisionChangeDigest string
		CycleStatus          string
		CycleEvidenceDigest  string
	}{
		RevisionChangeDigest: input.RevisionChangeDigest,
		CycleStatus:          input.Cycle.Status,
		CycleEvidenceDigest:  input.Cycle.EvidenceDigest,
	})
	if err != nil {
		output.MissingStage = "cycle-revision-change-digest"
		return output
	}
	candidate := GenerateJEVImprovementRevisionCandidate(JEVImprovementRevisionCandidateInput{
		Directive:           input.Directive,
		RevisionSource:      input.RevisionSource,
		RevisionChangeDigest: boundChangeDigest,
		NonAuthorizing:      true,
	})
	if err := candidate.Validate(); err != nil {
		output.MissingStage = "revision-candidate"
		return output
	}
	bindingDigest, err := Digest(struct {
		CycleEvidenceDigest     string
		CandidateDigest         string
		CandidateEvidenceDigest string
		BoundRevisionChangeDigest string
	}{
		CycleEvidenceDigest:       input.Cycle.EvidenceDigest,
		CandidateDigest:           candidate.CandidateDigest,
		CandidateEvidenceDigest:   candidate.EvidenceDigest,
		BoundRevisionChangeDigest: boundChangeDigest,
	})
	if err != nil {
		output.MissingStage = "cycle-revision-candidate-binding-digest"
		return output
	}
	output.Status = "bound"
	output.CycleStatus = input.Cycle.Status
	output.CycleEvidenceDigest = input.Cycle.EvidenceDigest
	output.CandidateStatus = candidate.Status
	output.CandidateDigest = candidate.CandidateDigest
	output.CandidateEvidenceDigest = candidate.EvidenceDigest
	output.BoundRevisionChangeDigest = boundChangeDigest
	output.BindingDigest = bindingDigest
	return output
}
