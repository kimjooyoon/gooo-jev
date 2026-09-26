package decision

import "fmt"

// ProposeImprovementCandidateFromDecisionConfidence bridges a reviewable
// confidence signal into the existing non-executing improvement pipeline.
func ProposeImprovementCandidateFromDecisionConfidence(
	signal DecisionConfidenceImprovementSignal,
	cycle ImprovementCycle,
	replay ImprovementCycleReplay,
	candidateReference string,
	sourceReference string,
) (ImprovementCandidate, error) {
	if err := signal.Validate(); err != nil {
		return ImprovementCandidate{}, fmt.Errorf("validate confidence improvement signal: %w", err)
	}
	if err := cycle.Validate(); err != nil {
		return ImprovementCandidate{}, fmt.Errorf("validate improvement cycle: %w", err)
	}
	if err := replay.Validate(); err != nil {
		return ImprovementCandidate{}, fmt.Errorf("validate improvement cycle replay: %w", err)
	}

	switch string(signal.Status) {
	case "candidate":
		candidate, err := ProposeImprovementCandidate(
			cycle,
			replay,
			candidateReference,
			sourceReference,
			CandidateReview,
		)
		if err != nil {
			return ImprovementCandidate{}, fmt.Errorf("propose improvement candidate: %w", err)
		}
		return candidate, nil
	case "no-change":
		return ImprovementCandidate{}, fmt.Errorf("confidence improvement signal has no candidate")
	case "unknown":
		return ImprovementCandidate{}, fmt.Errorf("confidence improvement signal is unknown")
	default:
		return ImprovementCandidate{}, fmt.Errorf("unsupported confidence improvement signal status %q", signal.Status)
	}
}
