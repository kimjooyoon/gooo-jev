package gooo

import "fmt"

// RevisionDirection is a bounded next-step hint, not an instruction to mutate source.
type RevisionDirection string

const (
	ClarifyRevision  RevisionDirection = "clarify"
	RepairRevision   RevisionDirection = "repair"
	PreserveRevision RevisionDirection = "preserve"
	RejectRevision   RevisionDirection = "reject"
)

// RevisionCandidate is a provenance-linked candidate for a later self-improvement step.
type RevisionCandidate struct {
	Status           string
	MissingStage     string
	Direction        RevisionDirection
	DecisionName     string
	DecisionID       string
	AssessmentDigest string
	EvidenceDigest   string
	CandidateDigest  string
	NonExecuting     bool
	NonAuthorizing   bool
}

// ProposeRevision derives a bounded candidate from a validated decision assessment.
func ProposeRevision(assessment DecisionAssessment, direction RevisionDirection) (RevisionCandidate, error) {
	candidate := RevisionCandidate{
		Status:           "UNKNOWN",
		MissingStage:     "revision-proposal",
		Direction:        direction,
		DecisionName:     assessment.DecisionName,
		DecisionID:       assessment.DecisionID,
		AssessmentDigest: assessment.AssessmentDigest,
		EvidenceDigest:   assessment.EvidenceDigest,
		NonExecuting:     true,
		NonAuthorizing:   true,
	}
	if err := assessment.Validate(); err != nil {
		candidate.MissingStage = "revision-assessment"
		return candidate, fmt.Errorf("gooo revision proposal: %w", err)
	}
	if !validRevisionDirection(direction) {
		candidate.MissingStage = "revision-direction"
		return candidate, fmt.Errorf("gooo revision proposal: direction %q is unsupported", direction)
	}

	candidate.Status = "BOUND"
	candidate.MissingStage = ""
	candidate.CandidateDigest = digestRevisionCandidate(candidate)
	return candidate, nil
}

// Validate confirms that a candidate is derived from a valid assessment and cannot authorize mutation.
func (c RevisionCandidate) Validate() error {
	if c.Status != "BOUND" {
		return fmt.Errorf("candidate status must be BOUND")
	}
	if c.MissingStage != "" {
		return fmt.Errorf("candidate missing stage must be empty")
	}
	if !validRevisionDirection(c.Direction) {
		return fmt.Errorf("candidate direction is invalid")
	}
	if c.DecisionName == "" || c.DecisionID == "" {
		return fmt.Errorf("candidate decision identity is incomplete")
	}
	if !validDigest(c.AssessmentDigest) || !validDigest(c.EvidenceDigest) || !validDigest(c.CandidateDigest) {
		return fmt.Errorf("candidate provenance digests are invalid")
	}
	if !c.NonExecuting || !c.NonAuthorizing {
		return fmt.Errorf("candidate must remain non-executing and non-authorizing")
	}
	if expected := digestRevisionCandidate(c); expected != c.CandidateDigest {
		return fmt.Errorf("candidate digest does not match its fields")
	}
	return nil
}

func validRevisionDirection(direction RevisionDirection) bool {
	switch direction {
	case ClarifyRevision, RepairRevision, PreserveRevision, RejectRevision:
		return true
	default:
		return false
	}
}

func digestRevisionCandidate(candidate RevisionCandidate) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s",
		candidate.Direction,
		candidate.DecisionName,
		candidate.DecisionID,
		candidate.AssessmentDigest,
		candidate.EvidenceDigest,
		candidate.NonExecuting,
	))
}
