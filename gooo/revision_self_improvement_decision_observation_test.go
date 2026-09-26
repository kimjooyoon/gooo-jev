package gooo

import "testing"

func selfImprovementDecisionInputs(t *testing.T) (RevisionSelfImprovementIteration, RevisionSelfImprovementOutcomeObservation, RevisionSelfImprovementFeedback) {
	t.Helper()
	iteration, application, metrics, generation := selfImprovementOutcomeInputs(t)
	outcome, err := ObserveRevisionSelfImprovementOutcome(iteration, application, metrics, generation)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementOutcome() error = %v", err)
	}
	_, _, feedback := selfImprovementIterationInputs(t)
	return iteration, outcome, feedback
}

func TestObserveRevisionSelfImprovementDecisionBindsNextAction(t *testing.T) {
	iteration, outcome, feedback := selfImprovementDecisionInputs(t)
	result, err := ObserveRevisionSelfImprovementDecision(iteration, outcome, feedback)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementDecision() error = %v", err)
	}
	if result.Status != "BOUND" ||
		result.DecisionSignal != feedback.FeedbackSignal ||
		result.DecisionReason != feedback.FeedbackReason ||
		result.OutcomeDigest != outcome.OutcomeDigest {
		t.Fatalf("unexpected self-improvement decision: %#v", result)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestObserveRevisionSelfImprovementDecisionRetainsOutcomeFailure(t *testing.T) {
	iteration, outcome, feedback := selfImprovementDecisionInputs(t)
	outcome.OutcomeDigest = digestString("tampered")
	result, err := ObserveRevisionSelfImprovementDecision(iteration, outcome, feedback)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementDecision() error = nil, want outcome failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-decision-outcome" {
		t.Fatalf("unexpected unknown outcome decision: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementDecisionRetainsFeedbackFailure(t *testing.T) {
	iteration, outcome, feedback := selfImprovementDecisionInputs(t)
	feedback.FeedbackDigest = digestString("tampered")
	result, err := ObserveRevisionSelfImprovementDecision(iteration, outcome, feedback)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementDecision() error = nil, want feedback failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-decision-feedback" {
		t.Fatalf("unexpected unknown feedback decision: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementDecisionRetainsIterationOutcomeLinkFailure(t *testing.T) {
	iteration, outcome, feedback := selfImprovementDecisionInputs(t)
	outcome.SourceDigest = digestString("other-source")
	outcome.OutcomeDigest = digestRevisionSelfImprovementOutcomeObservation(outcome)
	result, err := ObserveRevisionSelfImprovementDecision(iteration, outcome, feedback)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementDecision() error = nil, want iteration-outcome link failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-decision-iteration-outcome-link" {
		t.Fatalf("unexpected unknown iteration-outcome decision: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementDecisionRetainsIterationFeedbackLinkFailure(t *testing.T) {
	iteration, outcome, feedback := selfImprovementDecisionInputs(t)
	feedback.WindowDigest = digestString("other-window")
	feedback.FeedbackDigest = digestRevisionSelfImprovementFeedback(feedback)
	result, err := ObserveRevisionSelfImprovementDecision(iteration, outcome, feedback)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementDecision() error = nil, want iteration-feedback link failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-decision-iteration-feedback-link" {
		t.Fatalf("unexpected unknown iteration-feedback decision: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementDecisionIsDeterministic(t *testing.T) {
	iteration, outcome, feedback := selfImprovementDecisionInputs(t)
	first, err := ObserveRevisionSelfImprovementDecision(iteration, outcome, feedback)
	if err != nil {
		t.Fatalf("first ObserveRevisionSelfImprovementDecision() error = %v", err)
	}
	second, err := ObserveRevisionSelfImprovementDecision(iteration, outcome, feedback)
	if err != nil {
		t.Fatalf("second ObserveRevisionSelfImprovementDecision() error = %v", err)
	}
	if first.DecisionDigest != second.DecisionDigest {
		t.Fatal("same outcome and feedback produced different decision digest")
	}
}
