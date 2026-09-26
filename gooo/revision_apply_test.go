package gooo

import (
	"strings"
	"testing"
)

func revisionCandidateForApplication(t *testing.T) RevisionCandidate {
	t.Helper()
	return ProposeRevision(choiceAssessmentForRevision(t), ClarifyRevision)
}

func TestApplyRevisionBindsSourceEditAndReverseIR(t *testing.T) {
	edit := SourceEdit{
		Start:       Position{Line: 4, Column: 3},
		End:         Position{Line: 4, Column: 11},
		Replacement: "lineage",
	}
	edit.Digest = digestSourceEdit(edit)
	application, err := ApplyRevision(validContract, digestString(validContract), revisionCandidateForApplication(t), edit)
	if err != nil {
		t.Fatalf("ApplyRevision() error = %v", err)
	}
	if application.Status != "BOUND" || application.MissingStage != "" {
		t.Fatalf("unexpected application: %#v", application)
	}
	if application.SourceDigest == application.ProposedSourceDigest || application.InputIRDigest == application.ProposedIRDigest {
		t.Fatal("source and IR digests did not record the revision")
	}
	if !strings.Contains(application.ProposedSource, "# lineage are ignored") {
		t.Fatalf("proposed source missing edit: %q", application.ProposedSource)
	}
	if !application.NonExecuting || !application.NonAuthorizing {
		t.Fatalf("application must not authorize execution")
	}
	if err := application.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestApplyRevisionRejectsStaleSourcePrecondition(t *testing.T) {
	edit := SourceEdit{Start: Position{Line: 4, Column: 3}, End: Position{Line: 4, Column: 11}, Replacement: "lineage"}
	edit.Digest = digestSourceEdit(edit)
	application, err := ApplyRevision(validContract, digestString("stale-source"), revisionCandidateForApplication(t), edit)
	if err == nil {
		t.Fatal("ApplyRevision() error = nil, want precondition failure")
	}
	if application.Status != "UNKNOWN" || application.MissingStage != "revision-source-precondition" {
		t.Fatalf("unexpected precondition failure: %#v", application)
	}
}

func TestApplyRevisionRejectsInvalidRange(t *testing.T) {
	edit := SourceEdit{Start: Position{Line: 4, Column: 3}, End: Position{Line: 5, Column: 3}, Replacement: "lineage"}
	edit.Digest = digestSourceEdit(edit)
	application, err := ApplyRevision(validContract, digestString(validContract), revisionCandidateForApplication(t), edit)
	if err == nil {
		t.Fatal("ApplyRevision() error = nil, want range failure")
	}
	if application.MissingStage != "revision-range" {
		t.Fatalf("unexpected range failure: %#v", application)
	}
}

func TestApplyRevisionRetainsProposedParseUnknownStage(t *testing.T) {
	edit := SourceEdit{
		Start:       Position{Line: 5, Column: 8},
		End:         Position{Line: 5, Column: 20},
		Replacement: "MissingSpec",
	}
	edit.Digest = digestSourceEdit(edit)
	application, err := ApplyRevision(validContract, digestString(validContract), revisionCandidateForApplication(t), edit)
	if err == nil {
		t.Fatal("ApplyRevision() error = nil, want proposed parse failure")
	}
	if application.Status != "UNKNOWN" || application.MissingStage != "revision-reverse-observation" {
		t.Fatalf("unexpected reverse observation failure: %#v", application)
	}
	if application.ProposedSourceDigest == "" {
		t.Fatal("failed proposed source lost its digest")
	}
}

func TestValidateRejectsTamperedRevisionApplication(t *testing.T) {
	edit := SourceEdit{Start: Position{Line: 4, Column: 3}, End: Position{Line: 4, Column: 11}, Replacement: "lineage"}
	edit.Digest = digestSourceEdit(edit)
	application, err := ApplyRevision(validContract, digestString(validContract), revisionCandidateForApplication(t), edit)
	if err != nil {
		t.Fatalf("ApplyRevision() error = %v", err)
	}
	application.ProposedSourceDigest = digestString("tampered")
	if err := application.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want tamper rejection")
	}
}

func TestApplyRevisionIsDeterministic(t *testing.T) {
	edit := SourceEdit{Start: Position{Line: 4, Column: 3}, End: Position{Line: 4, Column: 11}, Replacement: "lineage"}
	edit.Digest = digestSourceEdit(edit)
	candidate := revisionCandidateForApplication(t)
	first, err := ApplyRevision(validContract, digestString(validContract), candidate, edit)
	if err != nil {
		t.Fatalf("first ApplyRevision() error = %v", err)
	}
	second, err := ApplyRevision(validContract, digestString(validContract), candidate, edit)
	if err != nil {
		t.Fatalf("second ApplyRevision() error = %v", err)
	}
	if first.ApplicationDigest != second.ApplicationDigest {
		t.Fatal("same source edit produced different application digest")
	}
}
