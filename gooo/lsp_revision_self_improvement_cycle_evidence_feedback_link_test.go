package gooo

import "testing"

func lspCycleEvidenceCoverageFeedbackInput(t *testing.T) RevisionSelfImprovementCycleEvidenceCoverageFeedbackObservation {
	t.Helper()
	bridge, coverage, feedback := cycleEvidenceCoverageFeedbackInputs(t)
	got, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedback(bridge, coverage, feedback)
	if err != nil {
		t.Fatalf("observe cycle evidence coverage feedback: %v", err)
	}
	return got
}

func TestObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackBindsProjection(t *testing.T) {
	link := lspCycleEvidenceCoverageFeedbackInput(t)
	got, err := ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedback(link)
	if err != nil {
		t.Fatalf("observe lsp cycle evidence coverage feedback: %v", err)
	}
	if got.Status != "BOUND" ||
		got.LinkSignal != "cycle-evidence-feedback-linked" ||
		!got.FeedbackAligned {
		t.Fatalf("unexpected lsp cycle evidence coverage feedback: %#v", got)
	}
	if got.BridgeDigest != link.BridgeDigest ||
		got.CoverageMetricDigest != link.CoverageMetricDigest ||
		got.FeedbackDigest != link.FeedbackDigest {
		t.Fatalf("projection lost stage digests: %#v", got)
	}
	if got.SourceDigest != link.SourceDigest ||
		got.CandidateSourceDigest != link.CandidateSourceDigest ||
		got.GeneratedIRDigest != link.GeneratedIRDigest ||
		got.ReverseObservationDigest != link.ReverseObservationDigest {
		t.Fatalf("projection lost evidence digests: %#v", got)
	}
	if got.FeedbackSignal != link.FeedbackSignal ||
		got.FeedbackReason != link.FeedbackReason ||
		!got.RequiresObservation ||
		!got.ReadOnly ||
		!got.NonExecuting ||
		!got.NonAuthorizing {
		t.Fatalf("projection lost feedback semantics: %#v", got)
	}
	if err := got.Validate(); err != nil {
		t.Fatalf("validate lsp cycle evidence coverage feedback: %v", err)
	}
}

func TestObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackPreservesUnknown(t *testing.T) {
	link := lspCycleEvidenceCoverageFeedbackInput(t)
	link.ReverseObservationDigest = digestString("tampered")
	got, err := ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedback(link)
	if err == nil {
		t.Fatal("expected tampered link error")
	}
	if got.Status != "UNKNOWN" ||
		got.MissingStage != "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-input" {
		t.Fatalf("unexpected unknown lsp cycle evidence coverage feedback: %#v", got)
	}
}

func TestObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackIsDeterministic(t *testing.T) {
	link := lspCycleEvidenceCoverageFeedbackInput(t)
	first, err := ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedback(link)
	if err != nil {
		t.Fatalf("first lsp cycle evidence coverage feedback: %v", err)
	}
	second, err := ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedback(link)
	if err != nil {
		t.Fatalf("second lsp cycle evidence coverage feedback: %v", err)
	}
	if first != second {
		t.Fatalf("lsp cycle evidence coverage feedback differs: %#v != %#v", first, second)
	}
}
