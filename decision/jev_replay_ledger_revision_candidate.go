package decision

import "strings"

// ExecutionEnvelopeJEVReplayLedgerRevisionCandidateInput binds a replay-ledger
// cycle observation to the existing non-executing revision candidate generator.
type ExecutionEnvelopeJEVReplayLedgerRevisionCandidateInput struct {
	Cycle                JEVReplayLedgerCycleObservation
	Directive            JEVImprovementDirectionDirective
	RevisionSource       string
	RevisionChangeDigest string
	NonAuthorizing       bool
}

// JEVReplayLedgerRevisionCandidateBinding preserves ledger and checkpoint
// evidence while exposing the existing candidate binding shape.
type JEVReplayLedgerRevisionCandidateBinding struct {
	Status                    string
	CycleStatus               string
	CycleEvidenceDigest       string
	CheckpointEvidenceDigest  string
	LedgerEvidenceDigest      string
	CandidateStatus           string
	CandidateDigest           string
	CandidateEvidenceDigest   string
	BoundRevisionChangeDigest string
	BindingDigest             string
	MissingStage              string
	NonExecuting              bool
	NonAuthorizing            bool
}

// GenerateJEVReplayLedgerRevisionCandidateFromCycle produces a candidate
// binding only for reviewable cycle outcomes, never an executable change.
func GenerateJEVReplayLedgerRevisionCandidateFromCycle(input ExecutionEnvelopeJEVReplayLedgerRevisionCandidateInput) JEVReplayLedgerRevisionCandidateBinding {
	output := JEVReplayLedgerRevisionCandidateBinding{
		Status:         "UNKNOWN",
		NonExecuting:   true,
		NonAuthorizing: true,
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
			output.MissingStage = "replay-ledger-cycle-observation"
		}
		return output
	}
	switch input.Cycle.Status {
	case "stable-for-review", "needs-revision":
		// Reviewable outcomes may produce evidence for a candidate.
	case "hold":
		output.MissingStage = "feedback-hold"
		return output
	default:
		output.MissingStage = "replay-ledger-cycle-status"
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
		RevisionChangeDigest      string
		CycleStatus               string
		CycleEvidenceDigest       string
		CheckpointEvidenceDigest  string
		LedgerEvidenceDigest      string
	}{
		RevisionChangeDigest:     input.RevisionChangeDigest,
		CycleStatus:              input.Cycle.Status,
		CycleEvidenceDigest:     input.Cycle.CycleEvidenceDigest,
		CheckpointEvidenceDigest: input.Cycle.CheckpointEvidenceDigest,
		LedgerEvidenceDigest:    input.Cycle.LedgerEvidenceDigest,
	})
	if err != nil {
		output.MissingStage = "replay-ledger-revision-change-digest"
		return output
	}
	candidate := GenerateJEVImprovementRevisionCandidate(JEVImprovementRevisionCandidateInput{
		Directive:            input.Directive,
		RevisionSource:       input.RevisionSource,
		RevisionChangeDigest: boundChangeDigest,
		NonAuthorizing:       true,
	})
	if err := candidate.Validate(); err != nil {
		output.MissingStage = "revision-candidate"
		return output
	}
	bindingDigest, err := Digest(struct {
		CycleEvidenceDigest       string
		CheckpointEvidenceDigest  string
		LedgerEvidenceDigest      string
		CandidateDigest           string
		CandidateEvidenceDigest   string
		BoundRevisionChangeDigest string
	}{
		CycleEvidenceDigest:       input.Cycle.CycleEvidenceDigest,
		CheckpointEvidenceDigest:  input.Cycle.CheckpointEvidenceDigest,
		LedgerEvidenceDigest:      input.Cycle.LedgerEvidenceDigest,
		CandidateDigest:           candidate.CandidateDigest,
		CandidateEvidenceDigest:   candidate.EvidenceDigest,
		BoundRevisionChangeDigest: boundChangeDigest,
	})
	if err != nil {
		output.MissingStage = "replay-ledger-revision-candidate-binding-digest"
		return output
	}
	output.Status = "bound"
	output.CycleStatus = input.Cycle.Status
	output.CycleEvidenceDigest = input.Cycle.CycleEvidenceDigest
	output.CheckpointEvidenceDigest = input.Cycle.CheckpointEvidenceDigest
	output.LedgerEvidenceDigest = input.Cycle.LedgerEvidenceDigest
	output.CandidateStatus = candidate.Status
	output.CandidateDigest = candidate.CandidateDigest
	output.CandidateEvidenceDigest = candidate.EvidenceDigest
	output.BoundRevisionChangeDigest = boundChangeDigest
	output.BindingDigest = bindingDigest
	return output
}