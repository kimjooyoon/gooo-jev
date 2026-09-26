package gooo

import "fmt"

type RevisionCandidateMaterialization struct {
	Status              string
	MissingStage        string
	SourceDigest        string
	FeedbackDigest      string
	DirectionDigest     string
	AssessmentDigest    string
	EvidenceDigest      string
	CandidateDigest     string
	Candidate           RevisionCandidate
	MaterializationDigest string
	NonExecuting        bool
	NonAuthorizing      bool
}

func MaterializeRevisionCandidate(feedback RevisionFeedback, proposal RevisionDirectionProposal, assessment DecisionAssessment) (RevisionCandidateMaterialization, error) {
	materialization := RevisionCandidateMaterialization{
		Status:           "UNKNOWN",
		MissingStage:     "revision-candidate-materialization",
		SourceDigest:     feedback.SourceDigest,
		FeedbackDigest:   feedback.FeedbackDigest,
		DirectionDigest:  proposal.DirectionDigest,
		AssessmentDigest: assessment.AssessmentDigest,
		EvidenceDigest:   assessment.EvidenceDigest,
		NonExecuting:     true,
		NonAuthorizing:   true,
	}
	if err := feedback.Validate(); err != nil {
		materialization.MissingStage = "revision-materialization-feedback"
		return materialization, fmt.Errorf("gooo candidate materialization: feedback: %w", err)
	}
	if err := proposal.Validate(); err != nil {
		materialization.MissingStage = "revision-materialization-direction"
		return materialization, fmt.Errorf("gooo candidate materialization: direction: %w", err)
	}
	if proposal.SourceDigest != feedback.SourceDigest || proposal.FeedbackDigest != feedback.FeedbackDigest {
		materialization.MissingStage = "revision-materialization-link"
		return materialization, fmt.Errorf("gooo candidate materialization: feedback and direction provenance do not match")
	}
	if err := assessment.Validate(); err != nil {
		materialization.MissingStage = "revision-materialization-assessment"
		return materialization, fmt.Errorf("gooo candidate materialization: assessment: %w", err)
	}
	candidate, err := ProposeRevision(assessment, proposal.Direction)
	if err != nil {
		materialization.MissingStage = "revision-materialization-candidate"
		return materialization, fmt.Errorf("gooo candidate materialization: candidate: %w", err)
	}
	materialization.Status = "BOUND"
	materialization.MissingStage = ""
	materialization.Candidate = candidate
	materialization.CandidateDigest = candidate.CandidateDigest
	materialization.MaterializationDigest = digestRevisionCandidateMaterialization(materialization)
	return materialization, nil
}

func (m RevisionCandidateMaterialization) Validate() error {
	if m.Status != "BOUND" {
		return fmt.Errorf("candidate materialization status must be BOUND")
	}
	if m.MissingStage != "" {
		return fmt.Errorf("candidate materialization missing stage must be empty")
	}
	if !m.NonExecuting || !m.NonAuthorizing {
		return fmt.Errorf("candidate materialization must remain non-executing and non-authorizing")
	}
	for name, digest := range map[string]string{
		"source": m.SourceDigest, "feedback": m.FeedbackDigest,
		"direction": m.DirectionDigest, "assessment": m.AssessmentDigest,
		"evidence": m.EvidenceDigest, "candidate": m.CandidateDigest,
		"materialization": m.MaterializationDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("candidate materialization %s digest is invalid", name)
		}
	}
	if err := m.Candidate.Validate(); err != nil {
		return fmt.Errorf("candidate materialization candidate: %w", err)
	}
	if m.Candidate.CandidateDigest != m.CandidateDigest {
		return fmt.Errorf("candidate materialization candidate digest is not linked")
	}
	if m.Candidate.AssessmentDigest != m.AssessmentDigest || m.Candidate.EvidenceDigest != m.EvidenceDigest {
		return fmt.Errorf("candidate materialization assessment evidence is not linked")
	}
	if digestRevisionCandidateMaterialization(m) != m.MaterializationDigest {
		return fmt.Errorf("candidate materialization digest does not match its fields")
	}
	return nil
}

func digestRevisionCandidateMaterialization(materialization RevisionCandidateMaterialization) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%t|%t",
		materialization.Status,
		materialization.MissingStage,
		materialization.SourceDigest,
		materialization.FeedbackDigest,
		materialization.DirectionDigest,
		materialization.AssessmentDigest,
		materialization.EvidenceDigest,
		materialization.CandidateDigest,
		materialization.NonExecuting,
		materialization.NonAuthorizing,
	))
}
