package decision

import "fmt"

// ProposeImprovementCandidateFromCounterexample reuses the existing generic
// non-executing candidate path for evidence that requires review.
func ProposeImprovementCandidateFromCounterexample(
	counterexample DecisionConfidenceImprovementCounterexample,
	cycle ImprovementCycle,
	replay ImprovementCycleReplay,
	candidateReference string,
	sourceReference string,
) (ImprovementCandidate, error) {
	if err := counterexample.Validate(); err != nil {
		return ImprovementCandidate{}, fmt.Errorf("validate confidence counterexample: %w", err)
	}
	candidate, err := ProposeImprovementCandidate(
		cycle,
		replay,
		candidateReference,
		sourceReference,
		CandidateReview,
	)
	if err != nil {
		return ImprovementCandidate{}, fmt.Errorf("propose review candidate from counterexample: %w", err)
	}
	return candidate, nil
}
