package gooo

import "testing"

func lspCycleFeedbackDecisionLinkInput(t *testing.T) RevisionSelfImprovementCycleFeedbackDecisionLinkObservation {
	t.Helper()
	bridge, decision, questionDigest := cycleFeedbackDecisionLinkInputs(t)
	link, err := ObserveRevisionSelfImprovementCycleFeedbackDecisionLink(bridge, decision, questionDigest)
	if err != nil {
		t.Fatalf("observe cycle feedback decision link: %v", err)
	}
	return link
}

func TestObserveLSPRevisionSelfImprovementCycleFeedbackDecisionLinkBindsProjection(t *testing.T) {
	link := lspCycleFeedbackDecisionLinkInput(t)
	got, err := ObserveLSPRevisionSelfImprovementCycleFeedbackDecisionLink(link)
	if err != nil {
		t.Fatalf("observe lsp cycle feedback decision link: %v", err)
	}
	if got.Status != "BOUND" || got.MissingStage != "" {
		t.Fatalf("unexpected lsp cycle feedback decision link: %#v", got)
	}
	if got.BridgeDigest != link.BridgeDigest || got.DecisionDigest != link.DecisionDigest || got.QuestionDigest != link.QuestionDigest {
		t.Fatalf("projection lost decision link digests: %#v", got)
	}
	if got.FeedbackSignal != "replan" || got.DecisionSignal != revisionSelfImprovementDecisionReceiptRequireReplan || !got.DecisionAligned {
		t.Fatalf("projection lost decision alignment: %#v", got)
	}
	if !got.ReadOnly || !got.NonExecuting || !got.NonAuthorizing {
		t.Fatalf("projection safety flags = %#v", got)
	}
	if err := got.Validate(); err != nil {
		t.Fatalf("validate lsp cycle feedback decision link: %v", err)
	}
}

func TestObserveLSPRevisionSelfImprovementCycleFeedbackDecisionLinkPreservesUnknown(t *testing.T) {
	link := lspCycleFeedbackDecisionLinkInput(t)
	link.DecisionDigest = digestString("tampered")
	got, err := ObserveLSPRevisionSelfImprovementCycleFeedbackDecisionLink(link)
	if err == nil {
		t.Fatal("expected tampered decision link error")
	}
	if got.Status != "UNKNOWN" || got.MissingStage != "lsp-revision-self-improvement-cycle-feedback-decision-link-input" {
		t.Fatalf("unexpected unknown lsp decision link: %#v", got)
	}
}

func TestObserveLSPRevisionSelfImprovementCycleFeedbackDecisionLinkIsDeterministic(t *testing.T) {
	link := lspCycleFeedbackDecisionLinkInput(t)
	first, err := ObserveLSPRevisionSelfImprovementCycleFeedbackDecisionLink(link)
	if err != nil {
		t.Fatalf("first lsp cycle feedback decision link: %v", err)
	}
	second, err := ObserveLSPRevisionSelfImprovementCycleFeedbackDecisionLink(link)
	if err != nil {
		t.Fatalf("second lsp cycle feedback decision link: %v", err)
	}
	if first != second {
		t.Fatalf("lsp cycle feedback decision links differ: %#v != %#v", first, second)
	}
}
