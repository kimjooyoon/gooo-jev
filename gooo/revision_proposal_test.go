package gooo

import "testing"

func choiceAssessmentForRevision(t *testing.T) DecisionAssessment {
	t.Helper()
	document, err := ParseDecision(validDecisionSource)
	if err != nil {
		t.Fatalf("ParseDecision() error = %v", err)
	}
	assessment, err := AssessDecision(document, DecisionObservation{
		DecisionName:   "ChooseReceipt",
		ObservedValue:  "DecisionReceipt",
		EvidenceDigest: digestString("revision-evidence"),
	})
	if err != nil {
		t.Fatalf("AssessDecision() error = %v", err)
	}
	return assessment
}

func TestProposeRevisionBindsAssessmentEvidence(t *testing.T) {
	candidate, err := ProposeRevision(choiceAssessmentForRevision(t), ClarifyRevision)
	if err != nil {
		t.Fatalf("ProposeRevision() error = %v", err)
	}
	if candidate.Status != "BOUND" || candidate.MissingStage != "" || candidate.Direction != ClarifyRevision {
		t.Fatalf("unexpected candidate: %#v", candidate)
	}
	if candidate.CandidateDigest == "" || !candidate.NonExecuting || !candidate.NonAuthorizing {
		t.Fatalf("missing safe candidate provenance: %#v", candidate)
	}
	if err := candidate.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestProposeRevisionSupportsBoundedDirections(t *testing.T) {
	assessment := choiceAssessmentForRevision(t)
	for _, direction := range []RevisionDirection{ClarifyRevision, RepairRevision, PreserveRevision, RejectRevision} {
		candidate, err := ProposeRevision(assessment, direction)
		if err != nil {
			t.Fatalf("ProposeRevision(%q) error = %v", direction, err)
		}
		if candidate.Direction != direction || candidate.CandidateDigest == "" {
			t.Fatalf("unexpected candidate for %q: %#v", direction, candidate)
		}
	}
}

func TestProposeRevisionRetainsUnknownDirectionStage(t *testing.T) {
	candidate, err := ProposeRevision(choiceAssessmentForRevision(t), RevisionDirection("mutate-now"))
	if err == nil {
		t.Fatal("ProposeRevision() error = nil, want direction rejection")
	}
	if candidate.Status != "UNKNOWN" || candidate.MissingStage != "revision-direction" {
		t.Fatalf("unexpected failure candidate: %#v", candidate)
	}
}

func TestProposeRevisionRetainsInvalidAssessmentStage(t *testing.T) {
	assessment := choiceAssessmentForRevision(t)
	assessment.EvidenceDigest = digestString("changed-evidence")
	candidate, err := ProposeRevision(assessment, RepairRevision)
	if err == nil {
		t.Fatal("ProposeRevision() error = nil, want assessment rejection")
	}
	if candidate.Status != "UNKNOWN" || candidate.MissingStage != "revision-assessment" {
		t.Fatalf("unexpected failure candidate: %#v", candidate)
	}
}

func TestValidateRejectsTamperedRevisionCandidate(t *testing.T) {
	candidate, err := ProposeRevision(choiceAssessmentForRevision(t), PreserveRevision)
	if err != nil {
		t.Fatalf("ProposeRevision() error = %v", err)
	}
	candidate.Direction = RejectRevision
	if err := candidate.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want tamper rejection")
	}
}

func TestProposeRevisionIsDeterministic(t *testing.T) {
	assessment := choiceAssessmentForRevision(t)
	first, err := ProposeRevision(assessment, RepairRevision)
	if err != nil {
		t.Fatalf("first ProposeRevision() error = %v", err)
	}
	second, err := ProposeRevision(assessment, RepairRevision)
	if err != nil {
		t.Fatalf("second ProposeRevision() error = %v", err)
	}
	if first.CandidateDigest != second.CandidateDigest {
		t.Fatalf("same assessment produced different candidate digest")
	}
}
