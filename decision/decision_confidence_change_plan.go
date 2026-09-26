package decision

import "fmt"

// DecisionConfidenceChangePlan describes a proposed write set without
// applying it. Source and write-set references remain independently auditable.
type DecisionConfidenceChangePlan struct {
	ProposalDigest     string
	SourceReference    string
	ChangeDigest       string
	WriteSetDigest     string
	NonExecuting       bool
	ChangePlanDigest   string
}

func PlanDecisionConfidenceChange(
	proposal DecisionConfidenceImprovementProposal,
	sourceReference string,
	writeSetDigest string,
) (DecisionConfidenceChangePlan, error) {
	if err := proposal.Validate(); err != nil {
		return DecisionConfidenceChangePlan{}, fmt.Errorf("validate improvement proposal: %w", err)
	}
	if sourceReference == "" || writeSetDigest == "" {
		return DecisionConfidenceChangePlan{}, fmt.Errorf("source and write-set references are required")
	}
	plan := DecisionConfidenceChangePlan{
		ProposalDigest:  proposal.ProposalDigest,
		SourceReference: sourceReference,
		ChangeDigest:    proposal.ProposedChangeDigest,
		WriteSetDigest:  writeSetDigest,
		NonExecuting:    true,
	}
	digest, err := Digest(plan)
	if err != nil {
		return DecisionConfidenceChangePlan{}, fmt.Errorf("digest change plan: %w", err)
	}
	plan.ChangePlanDigest = digest
	if err := plan.Validate(); err != nil {
		return DecisionConfidenceChangePlan{}, fmt.Errorf("validate change plan: %w", err)
	}
	return plan, nil
}

func (p DecisionConfidenceChangePlan) Validate() error {
	if p.ProposalDigest == "" || p.SourceReference == "" || p.ChangeDigest == "" || p.WriteSetDigest == "" {
		return fmt.Errorf("change plan references are required")
	}
	if !p.NonExecuting {
		return fmt.Errorf("change plan must be non-executing")
	}
	if p.ChangePlanDigest == "" {
		return fmt.Errorf("change plan digest is required")
	}
	withoutDigest := p
	withoutDigest.ChangePlanDigest = ""
	digest, err := Digest(withoutDigest)
	if err != nil {
		return fmt.Errorf("digest change plan: %w", err)
	}
	if digest != p.ChangePlanDigest {
		return fmt.Errorf("change plan digest mismatch")
	}
	return nil
}
