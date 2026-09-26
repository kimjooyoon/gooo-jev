package gooo

import "fmt"

type RevisionDirectionProposal struct {
	Status         string
	MissingStage   string
	SourceDigest   string
	FeedbackDigest string
	ChangeClass    string
	ActionHint     string
	Direction      RevisionDirection
	DirectionDigest string
	NonExecuting   bool
	NonAuthorizing bool
}

func ProposeRevisionDirection(feedback RevisionFeedback) (RevisionDirectionProposal, error) {
	proposal := RevisionDirectionProposal{
		Status:         "UNKNOWN",
		MissingStage:   "revision-direction-proposal",
		SourceDigest:   feedback.SourceDigest,
		FeedbackDigest: feedback.FeedbackDigest,
		ChangeClass:    feedback.ChangeClass,
		ActionHint:     feedback.ActionHint,
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if err := feedback.Validate(); err != nil {
		proposal.MissingStage = "revision-direction-feedback"
		return proposal, fmt.Errorf("gooo revision direction: feedback: %w", err)
	}
	proposal.Status = "BOUND"
	proposal.MissingStage = ""
	switch feedback.ActionHint {
	case "observe-more":
		proposal.Direction = ClarifyRevision
	case "inspect-ir":
		proposal.Direction = RepairRevision
	case "bound-scope":
		proposal.Direction = PreserveRevision
	default:
		proposal.Status = "UNKNOWN"
		proposal.MissingStage = "revision-direction-hint"
		return proposal, fmt.Errorf("gooo revision direction: unsupported feedback hint %q", feedback.ActionHint)
	}
	proposal.DirectionDigest = digestRevisionDirectionProposal(proposal)
	return proposal, nil
}

func (p RevisionDirectionProposal) Validate() error {
	if p.Status != "BOUND" {
		return fmt.Errorf("revision direction proposal status must be BOUND")
	}
	if p.MissingStage != "" {
		return fmt.Errorf("revision direction proposal missing stage must be empty")
	}
	if !p.NonExecuting || !p.NonAuthorizing {
		return fmt.Errorf("revision direction proposal must remain non-executing and non-authorizing")
	}
	if !validDigest(p.SourceDigest) || !validDigest(p.FeedbackDigest) || !validDigest(p.DirectionDigest) {
		return fmt.Errorf("revision direction proposal digests are invalid")
	}
	if p.ChangeClass != "localized" && p.ChangeClass != "structural" && p.ChangeClass != "broad" {
		return fmt.Errorf("revision direction proposal change class is invalid")
	}
	if p.ActionHint != "observe-more" && p.ActionHint != "inspect-ir" && p.ActionHint != "bound-scope" {
		return fmt.Errorf("revision direction proposal action hint is invalid")
	}
	if !validRevisionDirection(p.Direction) {
		return fmt.Errorf("revision direction proposal direction is invalid")
	}
	if digestRevisionDirectionProposal(p) != p.DirectionDigest {
		return fmt.Errorf("revision direction proposal digest does not match its fields")
	}
	return nil
}

func digestRevisionDirectionProposal(proposal RevisionDirectionProposal) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%t|%t",
		proposal.Status,
		proposal.MissingStage,
		proposal.SourceDigest,
		proposal.FeedbackDigest,
		proposal.ChangeClass,
		proposal.ActionHint,
		proposal.Direction,
		proposal.NonExecuting,
		proposal.NonAuthorizing,
	))
}
