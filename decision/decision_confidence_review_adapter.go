package decision

import "fmt"

// ReviewDecisionConfidenceCandidate reuses the generic non-executing review
// path after a confidence signal has produced a reviewable candidate.
func ReviewDecisionConfidenceCandidate(
	signal DecisionConfidenceImprovementSignal,
	cycle ImprovementCycle,
	replay ImprovementCycleReplay,
	candidateReference string,
	sourceReference string,
	reviewerReference string,
	reviewEvidenceDigest string,
	decision ImprovementReviewDecision,
) (ImprovementReviewReceipt, error) {
	candidate, err := ProposeImprovementCandidateFromDecisionConfidence(
		signal,
		cycle,
		replay,
		candidateReference,
		sourceReference,
	)
	if err != nil {
		return ImprovementReviewReceipt{}, fmt.Errorf("derive reviewable confidence candidate: %w", err)
	}

	receipt, err := ReviewImprovementCandidate(
		candidate,
		reviewerReference,
		reviewEvidenceDigest,
		decision,
	)
	if err != nil {
		return ImprovementReviewReceipt{}, fmt.Errorf("review confidence candidate: %w", err)
	}
	return receipt, nil
}
