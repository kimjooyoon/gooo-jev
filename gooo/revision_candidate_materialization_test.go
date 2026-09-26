package gooo

import "testing"

func materializationFeedbackAndProposal(t *testing.T) (RevisionFeedback, RevisionDirectionProposal) {
	t.Helper()
	feedback := revisionFeedbackForQuality(t)
	proposal, err := ProposeRevisionDirection(feedback)
	if err != nil {
		t.Fatalf("ProposeRevisionDirection() error = %v", err)
	}
	return feedback, proposal
}

func TestMaterializeRevisionCandidateBindsAllEvidence(t *testing.T) {
	feedback, proposal := materializationFeedbackAndProposal(t)
	materialization, err := MaterializeRevisionCandidate(feedback, proposal, choiceAssessmentForRevision(t))
	if err != nil {
		t.Fatalf("MaterializeRevisionCandidate() error = %v", err)
	}
	if materialization.Status != "BOUND" || materialization.Candidate.Status != "BOUND" {
		t.Fatalf("unexpected materialization: %#v", materialization)
	}
	if materialization.Candidate.Direction != RepairRevision {
		t.Fatalf("candidate direction = %q, want repair", materialization.Candidate.Direction)
	}
	if materialization.CandidateDigest != materialization.Candidate.CandidateDigest {
		t.Fatal("candidate digest was not linked")
	}
	if err := materialization.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestMaterializeRevisionCandidateRetainsUnknownAssessmentStage(t *testing.T) {
	feedback, proposal := materializationFeedbackAndProposal(t)
	assessment := choiceAssessmentForRevision(t)
	assessment.EvidenceDigest = digestString("tampered-evidence")
	materialization, err := MaterializeRevisionCandidate(feedback, proposal, assessment)
	if err == nil {
		t.Fatal("MaterializeRevisionCandidate() error = nil, want assessment failure")
	}
	if materialization.Status != "UNKNOWN" || materialization.MissingStage != "revision-materialization-assessment" {
		t.Fatalf("unexpected unknown materialization: %#v", materialization)
	}
}

func TestMaterializeRevisionCandidateIsDeterministic(t *testing.T) {
	feedback, proposal := materializationFeedbackAndProposal(t)
	assessment := choiceAssessmentForRevision(t)
	first, err := MaterializeRevisionCandidate(feedback, proposal, assessment)
	if err != nil {
		t.Fatalf("first MaterializeRevisionCandidate() error = %v", err)
	}
	second, err := MaterializeRevisionCandidate(feedback, proposal, assessment)
	if err != nil {
		t.Fatalf("second MaterializeRevisionCandidate() error = %v", err)
	}
	if first.MaterializationDigest != second.MaterializationDigest {
		t.Fatal("same evidence produced different materialization digest")
	}
}
