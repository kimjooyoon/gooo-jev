package gooo

import "testing"

func lineageRevisionInputs(t *testing.T) (LineageReceipt, RevisionCandidate) {
	t.Helper()
	document, assessment, candidate := lineageInputs(t)
	lineage, err := ObserveLineage(document, assessment, candidate)
	if err != nil {
		t.Fatalf("ObserveLineage() error = %v", err)
	}
	return lineage, candidate
}

func revisionCommentEdit() SourceEdit {
	edit := SourceEdit{
		Start:       Position{Line: 4, Column: 3},
		End:         Position{Line: 4, Column: 11},
		Replacement: "lineage",
	}
	edit.Digest = digestSourceEdit(edit)
	return edit
}

func TestApplyLineageRevisionBindsCompleteChain(t *testing.T) {
	lineage, candidate := lineageRevisionInputs(t)
	receipt, err := ApplyLineageRevision(validContract, lineage, candidate, revisionCommentEdit())
	if err != nil {
		t.Fatalf("ApplyLineageRevision() error = %v", err)
	}
	if receipt.Status != "BOUND" || receipt.MissingStage != "" {
		t.Fatalf("unexpected lineage revision: %#v", receipt)
	}
	if receipt.ResultDigest == "" || receipt.ApplicationDigest == "" {
		t.Fatalf("missing lineage revision evidence: %#v", receipt)
	}
	if err := receipt.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestApplyLineageRevisionRejectsDifferentSource(t *testing.T) {
	lineage, candidate := lineageRevisionInputs(t)
	receipt, err := ApplyLineageRevision(validContract+"\n", lineage, candidate, revisionCommentEdit())
	if err == nil {
		t.Fatal("ApplyLineageRevision() error = nil, want source binding failure")
	}
	if receipt.Status != "UNKNOWN" || receipt.MissingStage != "lineage-application-source" {
		t.Fatalf("unexpected source failure: %#v", receipt)
	}
}

func TestApplyLineageRevisionRejectsDifferentCandidate(t *testing.T) {
	lineage, candidate := lineageRevisionInputs(t)
	candidate.DecisionName = "DifferentDecision"
	candidate.CandidateDigest = digestRevisionCandidate(candidate)
	receipt, err := ApplyLineageRevision(validContract, lineage, candidate, revisionCommentEdit())
	if err == nil {
		t.Fatal("ApplyLineageRevision() error = nil, want candidate binding failure")
	}
	if receipt.MissingStage != "lineage-application-binding" {
		t.Fatalf("unexpected candidate failure: %#v", receipt)
	}
}

func TestApplyLineageRevisionRetainsLineageTamperStage(t *testing.T) {
	lineage, candidate := lineageRevisionInputs(t)
	lineage.ChainDigest = digestString("tampered-lineage")
	receipt, err := ApplyLineageRevision(validContract, lineage, candidate, revisionCommentEdit())
	if err == nil {
		t.Fatal("ApplyLineageRevision() error = nil, want lineage failure")
	}
	if receipt.Status != "UNKNOWN" || receipt.MissingStage != "lineage-application-lineage" {
		t.Fatalf("unexpected lineage failure: %#v", receipt)
	}
}

func TestApplyLineageRevisionRetainsProposedUnknownStage(t *testing.T) {
	lineage, candidate := lineageRevisionInputs(t)
	edit := SourceEdit{
		Start:       Position{Line: 5, Column: 8},
		End:         Position{Line: 5, Column: 20},
		Replacement: "MissingSpec",
	}
	edit.Digest = digestSourceEdit(edit)
	receipt, err := ApplyLineageRevision(validContract, lineage, candidate, edit)
	if err == nil {
		t.Fatal("ApplyLineageRevision() error = nil, want proposed source failure")
	}
	if receipt.MissingStage != "revision-reverse-observation" || receipt.ProposedSourceDigest == "" {
		t.Fatalf("unexpected proposed source failure: %#v", receipt)
	}
}

func TestValidateRejectsTamperedLineageRevision(t *testing.T) {
	lineage, candidate := lineageRevisionInputs(t)
	receipt, err := ApplyLineageRevision(validContract, lineage, candidate, revisionCommentEdit())
	if err != nil {
		t.Fatalf("ApplyLineageRevision() error = %v", err)
	}
	receipt.ProposedSource = "tampered"
	if err := receipt.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want tamper rejection")
	}
}

func TestApplyLineageRevisionIsDeterministic(t *testing.T) {
	lineage, candidate := lineageRevisionInputs(t)
	first, err := ApplyLineageRevision(validContract, lineage, candidate, revisionCommentEdit())
	if err != nil {
		t.Fatalf("first ApplyLineageRevision() error = %v", err)
	}
	second, err := ApplyLineageRevision(validContract, lineage, candidate, revisionCommentEdit())
	if err != nil {
		t.Fatalf("second ApplyLineageRevision() error = %v", err)
	}
	if first.ResultDigest != second.ResultDigest {
		t.Fatal("same complete lineage produced different result digest")
	}
}
