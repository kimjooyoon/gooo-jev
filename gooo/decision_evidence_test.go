package gooo

import "testing"

func TestAssessDecisionBindsChoiceEvidence(t *testing.T) {
	document, err := ParseDecision(validDecisionSource)
	if err != nil {
		t.Fatalf("ParseDecision() error = %v", err)
	}
	evidence := digestString("choice-evidence")
	assessment, err := AssessDecision(document, DecisionObservation{
		DecisionName:   "ChooseReceipt",
		ObservedValue:  "DecisionReceipt",
		EvidenceDigest: evidence,
	})
	if err != nil {
		t.Fatalf("AssessDecision() error = %v", err)
	}
	if assessment.Status != "BOUND" || assessment.MissingStage != "" || assessment.DecisionKind != ChoiceDecision {
		t.Fatalf("unexpected assessment: %#v", assessment)
	}
	if assessment.AssessmentDigest == "" || !assessment.NonExecuting || !assessment.NonAuthorizing {
		t.Fatalf("missing safe assessment evidence: %#v", assessment)
	}
	if err := assessment.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestAssessDecisionBindsScoreAndBooleanEvidence(t *testing.T) {
	source := "package jevdecision\n" +
		"namespace jevdecision\n" +
		"entity DecisionSpec id \"gooo://jev/decision/spec\"\n" +
		"entity DecisionReceipt id \"gooo://jev/decision/receipt\"\n" +
		"decision Confidence kind score id \"gooo://jev/decision/score\"\n" +
		"decision Approved kind boolean id \"gooo://jev/decision/boolean\"\n" +
		"activity ObserveDecision(DecisionSpec) -> DecisionReceipt\n"
	document, err := ParseDecision(source)
	if err != nil {
		t.Fatalf("ParseDecision() error = %v", err)
	}
	score := 0.75
	scoreAssessment, err := AssessDecision(document, DecisionObservation{
		DecisionName:   "Confidence",
		Score:          &score,
		EvidenceDigest: digestString("score-evidence"),
	})
	if err != nil {
		t.Fatalf("score AssessDecision() error = %v", err)
	}
	if !scoreAssessment.HasScore || scoreAssessment.Score != score {
		t.Fatalf("unexpected score assessment: %#v", scoreAssessment)
	}
	approved := false
	booleanAssessment, err := AssessDecision(document, DecisionObservation{
		DecisionName:   "Approved",
		Boolean:        &approved,
		EvidenceDigest: digestString("boolean-evidence"),
	})
	if err != nil {
		t.Fatalf("boolean AssessDecision() error = %v", err)
	}
	if !booleanAssessment.HasBoolean || booleanAssessment.Boolean {
		t.Fatalf("unexpected boolean assessment: %#v", booleanAssessment)
	}
}

func TestAssessDecisionRetainsMissingEvidenceStage(t *testing.T) {
	document, err := ParseDecision(validDecisionSource)
	if err != nil {
		t.Fatalf("ParseDecision() error = %v", err)
	}
	assessment, err := AssessDecision(document, DecisionObservation{
		DecisionName:  "ChooseReceipt",
		ObservedValue: "DecisionReceipt",
	})
	if err == nil {
		t.Fatal("AssessDecision() error = nil, want missing evidence")
	}
	if assessment.Status != "UNKNOWN" || assessment.MissingStage != "decision-evidence" {
		t.Fatalf("unexpected failure assessment: %#v", assessment)
	}
}

func TestAssessDecisionRejectsInvalidScore(t *testing.T) {
	source := "package jevdecision\n" +
		"namespace jevdecision\n" +
		"entity DecisionSpec id \"gooo://jev/decision/spec\"\n" +
		"entity DecisionReceipt id \"gooo://jev/decision/receipt\"\n" +
		"decision Confidence kind score id \"gooo://jev/decision/score\"\n" +
		"activity ObserveDecision(DecisionSpec) -> DecisionReceipt\n"
	document, err := ParseDecision(source)
	if err != nil {
		t.Fatalf("ParseDecision() error = %v", err)
	}
	score := 1.5
	assessment, err := AssessDecision(document, DecisionObservation{
		DecisionName:   "Confidence",
		Score:          &score,
		EvidenceDigest: digestString("score-evidence"),
	})
	if err == nil {
		t.Fatal("AssessDecision() error = nil, want score rejection")
	}
	if assessment.MissingStage != "decision-score" {
		t.Fatalf("unexpected score failure: %#v", assessment)
	}
}

func TestAssessDecisionRejectsMismatchedPayload(t *testing.T) {
	document, err := ParseDecision(validDecisionSource)
	if err != nil {
		t.Fatalf("ParseDecision() error = %v", err)
	}
	score := 0.5
	assessment, err := AssessDecision(document, DecisionObservation{
		DecisionName:   "ChooseReceipt",
		ObservedValue:  "DecisionReceipt",
		Score:          &score,
		EvidenceDigest: digestString("shape-evidence"),
	})
	if err == nil {
		t.Fatal("AssessDecision() error = nil, want shape rejection")
	}
	if assessment.MissingStage != "decision-shape" {
		t.Fatalf("unexpected shape failure: %#v", assessment)
	}
}

func TestAssessDecisionRejectsUnknownDecision(t *testing.T) {
	document, err := ParseDecision(validDecisionSource)
	if err != nil {
		t.Fatalf("ParseDecision() error = %v", err)
	}
	assessment, err := AssessDecision(document, DecisionObservation{
		DecisionName:   "Missing",
		ObservedValue:  "DecisionReceipt",
		EvidenceDigest: digestString("unknown-evidence"),
	})
	if err == nil {
		t.Fatal("AssessDecision() error = nil, want resolution failure")
	}
	if assessment.MissingStage != "decision-resolution" {
		t.Fatalf("unexpected resolution failure: %#v", assessment)
	}
}

func TestValidateRejectsTamperedAssessment(t *testing.T) {
	document, err := ParseDecision(validDecisionSource)
	if err != nil {
		t.Fatalf("ParseDecision() error = %v", err)
	}
	assessment, err := AssessDecision(document, DecisionObservation{
		DecisionName:   "ChooseReceipt",
		ObservedValue:  "DecisionReceipt",
		EvidenceDigest: digestString("tamper-evidence"),
	})
	if err != nil {
		t.Fatalf("AssessDecision() error = %v", err)
	}
	assessment.EvidenceDigest = digestString("different-evidence")
	if err := assessment.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want tamper rejection")
	}
}
