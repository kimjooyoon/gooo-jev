package gooo

import "testing"

const lineageRevisionSource = "package jevdecision\n" +
	"namespace jevdecision\n" +
	"# lineage source\n" +
	"entity DecisionSpec id \"gooo://jev/decision/spec\"\n" +
	"entity DecisionReceipt id \"gooo://jev/decision/receipt\"\n" +
	"decision ChooseReceipt kind choice id \"gooo://jev/decision/choose-receipt\"\n" +
	"activity ObserveDecision(DecisionSpec) -> DecisionReceipt\n"

func lineageRevisionInputs(t *testing.T) (LineageReceipt, RevisionCandidate) {
	t.Helper()
	document, err := ParseDecision(lineageRevisionSource)
	if err != nil {
		t.Fatalf("ParseDecision() error = %v", err)
	}
	assessment, err := AssessDecision(document, DecisionObservation{
		DecisionName:   "ChooseReceipt",
		ObservedValue:  "DecisionReceipt",
		EvidenceDigest: digestString("lineage-revision-evidence"),
	})
	if err != nil {
		t.Fatalf("AssessDecision() error = %v", err)
	}
	candidate, err := ProposeRevision(assessment, RepairRevision)
	if err != nil {
		t.Fatalf("ProposeRevision() error = %v", err)
	}
	lineage, err := ObserveLineage(document, assessment, candidate)
	if err != nil {
		t.Fatalf("ObserveLineage() error = %v", err)
	}
	return lineage, candidate
}

func revisionCommentEdit() SourceEdit {
	edit := SourceEdit{
		Start:       Position{Line: 3, Column: 3},
		End:         Position{Line: 3, Column: 10},
		Replacement: "origin",
	}
	edit.Digest = digestSourceEdit(edit)
	return edit
}

func TestApplyLineageRevisionBindsCompleteChain(t *testing.T) {
	lineage, candidate := lineageRevisionInputs(t)
	receipt, err := ApplyLineageRevision(lineageRevisionSource, lineage, candidate, revisionCommentEdit())
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
	receipt, err := ApplyLineageRevision(lineageRevisionSource+"\n", lineage, candidate, revisionCommentEdit())
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
	receipt, err := ApplyLineageRevision(lineageRevisionSource, lineage, candidate, revisionCommentEdit())
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
	receipt, err := ApplyLineageRevision(lineageRevisionSource, lineage, candidate, revisionCommentEdit())
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
		Start:       Position{Line: 4, Column: 8},
		End:         Position{Line: 4, Column: 20},
		Replacement: "MissingSpec",
	}
	edit.Digest = digestSourceEdit(edit)
	receipt, err := ApplyLineageRevision(lineageRevisionSource, lineage, candidate, edit)
	if err == nil {
		t.Fatal("ApplyLineageRevision() error = nil, want proposed source failure")
	}
	if receipt.MissingStage != "revision-reverse-observation" || receipt.ProposedSourceDigest == "" {
		t.Fatalf("unexpected proposed source failure: %#v", receipt)
	}
}

func TestValidateRejectsTamperedLineageRevision(t *testing.T) {
	lineage, candidate := lineageRevisionInputs(t)
	receipt, err := ApplyLineageRevision(lineageRevisionSource, lineage, candidate, revisionCommentEdit())
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
	first, err := ApplyLineageRevision(lineageRevisionSource, lineage, candidate, revisionCommentEdit())
	if err != nil {
		t.Fatalf("first ApplyLineageRevision() error = %v", err)
	}
	second, err := ApplyLineageRevision(lineageRevisionSource, lineage, candidate, revisionCommentEdit())
	if err != nil {
		t.Fatalf("second ApplyLineageRevision() error = %v", err)
	}
	if first.ResultDigest != second.ResultDigest {
		t.Fatal("same complete lineage produced different result digest")
	}
}
