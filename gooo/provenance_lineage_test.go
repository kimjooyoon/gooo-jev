package gooo

import "testing"

func lineageInputs(t *testing.T) (DecisionDocumentIR, DecisionAssessment, RevisionCandidate) {
	t.Helper()
	document, err := ParseDecision(validDecisionSource)
	if err != nil {
		t.Fatalf("ParseDecision() error = %v", err)
	}
	assessment, err := AssessDecision(document, DecisionObservation{
		DecisionName:   "ChooseReceipt",
		ObservedValue:  "DecisionReceipt",
		EvidenceDigest: digestString("lineage-evidence"),
	})
	if err != nil {
		t.Fatalf("AssessDecision() error = %v", err)
	}
	candidate, err := ProposeRevision(assessment, RepairRevision)
	if err != nil {
		t.Fatalf("ProposeRevision() error = %v", err)
	}
	return document, assessment, candidate
}

func TestObserveLineageBindsAllEvidenceStages(t *testing.T) {
	document, assessment, candidate := lineageInputs(t)
	receipt, err := ObserveLineage(document, assessment, candidate)
	if err != nil {
		t.Fatalf("ObserveLineage() error = %v", err)
	}
	if receipt.Status != "BOUND" || receipt.MissingStage != "" || len(receipt.Edges) != 4 {
		t.Fatalf("unexpected lineage receipt: %#v", receipt)
	}
	if receipt.ChainDigest == "" || !receipt.NonExecuting || !receipt.NonAuthorizing {
		t.Fatalf("missing lineage safety evidence: %#v", receipt)
	}
	if err := receipt.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestObserveLineageRetainsBindingStage(t *testing.T) {
	document, assessment, candidate := lineageInputs(t)
	candidate.DecisionName = "DifferentDecision"
	candidate.CandidateDigest = digestRevisionCandidate(candidate)
	receipt, err := ObserveLineage(document, assessment, candidate)
	if err == nil {
		t.Fatal("ObserveLineage() error = nil, want binding failure")
	}
	if receipt.Status != "UNKNOWN" || receipt.MissingStage != "lineage-binding" {
		t.Fatalf("unexpected binding failure: %#v", receipt)
	}
}

func TestValidateRejectsTamperedLineageEdge(t *testing.T) {
	document, assessment, candidate := lineageInputs(t)
	receipt, err := ObserveLineage(document, assessment, candidate)
	if err != nil {
		t.Fatalf("ObserveLineage() error = %v", err)
	}
	receipt.Edges[0].ToDigest = digestString("tampered-edge")
	if err := receipt.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want edge tamper rejection")
	}
}

func TestObserveLineageIsDeterministic(t *testing.T) {
	document, assessment, candidate := lineageInputs(t)
	first, err := ObserveLineage(document, assessment, candidate)
	if err != nil {
		t.Fatalf("first ObserveLineage() error = %v", err)
	}
	second, err := ObserveLineage(document, assessment, candidate)
	if err != nil {
		t.Fatalf("second ObserveLineage() error = %v", err)
	}
	if first.ChainDigest != second.ChainDigest {
		t.Fatal("same evidence chain produced different digest")
	}
}
