package decision

import "fmt"

// ReviewDecisionConfidenceCounterexample reuses the generic review receipt
// after a counterexample has become a non-executing candidate.
func ReviewDecisionConfidenceCounterexample(
	counterexample DecisionConfidenceImprovementCounterexample,
	cycle ImprovementCycle,
	replay ImprovementCycleReplay,
	candidateReference string,
	sourceReference string,
	reviewerReference string,
	reviewEvidenceDigest string,
	decision ImprovementReviewDecision,
) (ImprovementReviewReceipt, error) {
	candidate, err := ProposeImprovementCandidateFromCounterexample(
		counterexample,
		cycle,
		replay,
		candidateReference,
		sourceReference,
	)
	if err != nil {
		return ImprovementReviewReceipt{}, fmt.Errorf("derive counterexample review candidate: %w", err)
	}
	receipt, err := ReviewImprovementCandidate(
		candidate,
		reviewerReference,
		reviewEvidenceDigest,
		decision,
	)
	if err != nil {
		return ImprovementReviewReceipt{}, fmt.Errorf("review counterexample candidate: %w", err)
	}
	return receipt, nil
}
