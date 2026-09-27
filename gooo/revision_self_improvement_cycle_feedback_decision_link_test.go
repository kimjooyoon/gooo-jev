package gooo

import (
	"strings"
	"testing"
)

func cycleFeedbackDecisionLinkInputs(t *testing.T) (
	RevisionSelfImprovementCycleFeedbackBridgeObservation,
	RevisionSelfImprovementDecisionReceiptObservation,
	string,
) {
	t.Helper()
	cycle := cycleMetricCycle(t)
	metricInput := cycleMetricInput(t)
	metricInput.Direction = "lower-is-better"
	metricInput.BaselineValue = 12
	metricInput.CandidateValue = 15
	metric, err := ObserveRevisionSelfImprovementCycleMetric(cycle, metricInput)
	if err != nil {
		t.Fatalf("observe regressed cycle metric: %v", err)
	}
	feedback, err := ObserveRevisionSelfImprovementCycleMetricFeedback(metric)
	if err != nil {
		t.Fatalf("observe regressed cycle feedback: %v", err)
	}
	bridge, err := ObserveRevisionSelfImprovementCycleFeedbackBridge(cycle, metric, feedback)
	if err != nil {
		t.Fatalf("observe cycle feedback bridge: %v", err)
	}
	replan, nextIteration := provenanceNextIterationInputs(t)
	nextBoundary, err := ObserveRevisionSelfImprovementProvenanceNextIteration(replan, nextIteration)
	if err != nil {
		t.Fatalf("observe provenance next iteration: %v", err)
	}
	questionDigest := digestString(strings.Join([]string{"cycle-feedback", bridge.FeedbackDigest}, "|"))
	decision, err := ObserveRevisionSelfImprovementProvenanceDecisionReceipt(
		nextBoundary,
		questionDigest,
		revisionSelfImprovementDecisionReceiptRequireReplan,
		"cycle metric regression requires replan",
	)
	if err != nil {
		t.Fatalf("observe decision receipt: %v", err)
	}
	return bridge, decision, questionDigest
}

func TestObserveRevisionSelfImprovementCycleFeedbackDecisionLinkBindsReplan(t *testing.T) {
	bridge, decision, questionDigest := cycleFeedbackDecisionLinkInputs(t)
	got, err := ObserveRevisionSelfImprovementCycleFeedbackDecisionLink(bridge, decision, questionDigest)
	if err != nil {
		t.Fatalf("observe cycle feedback decision link: %v", err)
	}
	if got.Status != "BOUND" || got.LinkSignal != "cycle-feedback-decision-linked" || !got.DecisionAligned {
		t.Fatalf("unexpected cycle feedback decision link: %#v", got)
	}
	if got.BridgeDigest != bridge.ObservationDigest || got.DecisionDigest != decision.ObservationDigest || got.QuestionDigest != questionDigest {
		t.Fatalf("link lost provenance digests: %#v", got)
	}
	if got.FeedbackSignal != "replan" || got.DecisionSignal != revisionSelfImprovementDecisionReceiptRequireReplan {
		t.Fatalf("replan signal was not preserved: %#v", got)
	}
	if !got.NonExecuting || !got.NonAuthorizing {
		t.Fatalf("link safety flags = %#v", got)
	}
	if err := got.Validate(); err != nil {
		t.Fatalf("validate cycle feedback decision link: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCycleFeedbackDecisionLinkPreservesUnknown(t *testing.T) {
	bridge, decision, questionDigest := cycleFeedbackDecisionLinkInputs(t)
	decision.DecisionDigest = digestString("tampered")
	got, err := ObserveRevisionSelfImprovementCycleFeedbackDecisionLink(bridge, decision, questionDigest)
	if err == nil {
		t.Fatal("expected tampered decision receipt error")
	}
	if got.Status != "UNKNOWN" || got.MissingStage != "revision-self-improvement-cycle-feedback-decision-link-decision" {
		t.Fatalf("unexpected unknown cycle feedback decision link: %#v", got)
	}
}

func TestObserveRevisionSelfImprovementCycleFeedbackDecisionLinkIsDeterministic(t *testing.T) {
	bridge, decision, questionDigest := cycleFeedbackDecisionLinkInputs(t)
	first, err := ObserveRevisionSelfImprovementCycleFeedbackDecisionLink(bridge, decision, questionDigest)
	if err != nil {
		t.Fatalf("first cycle feedback decision link: %v", err)
	}
	second, err := ObserveRevisionSelfImprovementCycleFeedbackDecisionLink(bridge, decision, questionDigest)
	if err != nil {
		t.Fatalf("second cycle feedback decision link: %v", err)
	}
	if first != second {
		t.Fatalf("cycle feedback decision links differ: %#v != %#v", first, second)
	}
}
