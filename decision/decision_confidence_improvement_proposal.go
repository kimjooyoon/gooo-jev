package decision

import "fmt"

// DecisionConfidenceImprovementProposal is immutable proposal evidence. It
// names a proposed change but never applies it or grants execution.
type DecisionConfidenceImprovementProposal struct {
	CounterexampleDigest string
	CandidateDigest     string
	ReviewDigest        string
	ProposedChangeDigest string
	NonExecuting        bool
	ProposalDigest      string
}

func ProposeDecisionConfidenceImprovement(
	counterexample DecisionConfidenceImprovementCounterexample,
	candidate ImprovementCandidate,
	review ImprovementReviewReceipt,
	proposedChangeDigest string,
) (DecisionConfidenceImprovementProposal, error) {
	if err := counterexample.Validate(); err != nil {
		return DecisionConfidenceImprovementProposal{}, fmt.Errorf("validate counterexample: %w", err)
	}
	if err := candidate.Validate(); err != nil {
		return DecisionConfidenceImprovementProposal{}, fmt.Errorf("validate improvement candidate: %w", err)
	}
	if err := review.Validate(); err != nil {
		return DecisionConfidenceImprovementProposal{}, fmt.Errorf("validate improvement review: %w", err)
	}
	if review.CandidateDigest != candidate.CandidateDigest {
		return DecisionConfidenceImprovementProposal{}, fmt.Errorf("review does not reference candidate")
	}
	if review.Decision != ReviewPassed {
		return DecisionConfidenceImprovementProposal{}, fmt.Errorf("improvement review is not passed")
	}
	if review.ExecutionGranted || proposedChangeDigest == "" {
		return DecisionConfidenceImprovementProposal{}, fmt.Errorf("proposal must remain non-executing and include a change digest")
	}
	proposal := DecisionConfidenceImprovementProposal{
		CounterexampleDigest:  counterexample.CounterexampleDigest,
		CandidateDigest:      candidate.CandidateDigest,
		ReviewDigest:         review.ReviewDigest,
		ProposedChangeDigest: proposedChangeDigest,
		NonExecuting:         true,
	}
	digest, err := Digest(proposal)
	if err != nil {
		return DecisionConfidenceImprovementProposal{}, fmt.Errorf("digest improvement proposal: %w", err)
	}
	proposal.ProposalDigest = digest
	if err := proposal.Validate(); err != nil {
		return DecisionConfidenceImprovementProposal{}, fmt.Errorf("validate improvement proposal: %w", err)
	}
	return proposal, nil
}

func (p DecisionConfidenceImprovementProposal) Validate() error {
	if p.CounterexampleDigest == "" || p.CandidateDigest == "" || p.ReviewDigest == "" || p.ProposedChangeDigest == "" {
		return fmt.Errorf("proposal digests are required")
	}
	if !p.NonExecuting {
		return fmt.Errorf("improvement proposal must be non-executing")
	}
	if p.ProposalDigest == "" {
		return fmt.Errorf("proposal digest is required")
	}
	withoutDigest := p
	withoutDigest.ProposalDigest = ""
	digest, err := Digest(withoutDigest)
	if err != nil {
		return fmt.Errorf("digest improvement proposal: %w", err)
	}
	if digest != p.ProposalDigest {
		return fmt.Errorf("proposal digest mismatch")
	}
	return nil
}
