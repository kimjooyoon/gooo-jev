package gooo

import "testing"

func lspCycleEvidenceCoverageFeedbackDecisionBoundaryInput(t *testing.T) RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryObservation {
	t.Helper()
	cycle, evidenceFeedback, boundary := cycleEvidenceCoverageFeedbackDecisionBoundaryInputs(t)
	got, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundary(cycle, evidenceFeedback, boundary)
	if err != nil {
		t.Fatalf("observe cycle evidence feedback decision boundary: %v", err)
	}
	return got
}

func TestObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryBindsProjection(t *testing.T) {
	boundary := lspCycleEvidenceCoverageFeedbackDecisionBoundaryInput(t)
	got, err := ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundary(boundary)
	if err != nil {
		t.Fatalf("observe lsp cycle evidence feedback decision boundary: %v", err)
	}
	if got.Status != "BOUND" ||
		got.BoundarySignal != "cycle-evidence-feedback-decision-boundary-linked" ||
		!got.EvidenceFeedbackAligned ||
		!got.BoundaryAligned {
		t.Fatalf("unexpected lsp cycle evidence feedback decision boundary: %#v", got)
	}
	if got.CycleDigest != boundary.CycleDigest ||
		got.BridgeDigest != boundary.BridgeDigest ||
		got.EvidenceFeedbackDigest != boundary.EvidenceFeedbackDigest ||
		got.DecisionBoundaryDigest != boundary.DecisionBoundaryDigest {
		t.Fatalf("projection lost provenance links: %#v", got)
	}
	if got.CoverageSignal != "evidence-complete" ||
		got.FeedbackSignal != boundary.FeedbackSignal ||
		got.DecisionSignal != boundary.DecisionSignal ||
		!got.ReadOnly ||
		!got.NonExecuting ||
		!got.NonAuthorizing {
		t.Fatalf("projection lost boundary semantics: %#v", got)
	}
	if err := got.Validate(); err != nil {
		t.Fatalf("validate lsp cycle evidence feedback decision boundary: %v", err)
	}
}

func TestObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryPreservesUnknown(t *testing.T) {
	boundary := lspCycleEvidenceCoverageFeedbackDecisionBoundaryInput(t)
	boundary.DecisionBoundaryDigest = digestString("tampered")
	got, err := ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundary(boundary)
	if err == nil {
		t.Fatal("expected tampered boundary error")
	}
	if got.Status != "UNKNOWN" ||
		got.MissingStage != "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-input" {
		t.Fatalf("unexpected unknown lsp cycle evidence feedback decision boundary: %#v", got)
	}
}

func TestObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryIsDeterministic(t *testing.T) {
	boundary := lspCycleEvidenceCoverageFeedbackDecisionBoundaryInput(t)
	first, err := ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundary(boundary)
	if err != nil {
		t.Fatalf("first lsp cycle evidence feedback decision boundary: %v", err)
	}
	second, err := ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundary(boundary)
	if err != nil {
		t.Fatalf("second lsp cycle evidence feedback decision boundary: %v", err)
	}
	if first != second {
		t.Fatalf("lsp cycle evidence feedback decision boundaries differ: %#v != %#v", first, second)
	}
}
