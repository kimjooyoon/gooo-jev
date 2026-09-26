package gooo

import "testing"

func selfImprovementIterationInputs(t *testing.T) (string, RevisionSelfImprovementHistory, RevisionSelfImprovementFeedback) {
	t.Helper()
	windows, feedback := selfImprovementHistoryInputs(t)
	history, err := ObserveRevisionSelfImprovementHistory(windows, feedback)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementHistory() error = %v", err)
	}
	return validContract, history, feedback[len(feedback)-1]
}

func TestObserveRevisionSelfImprovementIterationBindsCurrentSnapshot(t *testing.T) {
	source, history, feedback := selfImprovementIterationInputs(t)
	iteration, err := ObserveRevisionSelfImprovementIteration(source, history, feedback)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementIteration() error = %v", err)
	}
	if iteration.Status != "BOUND" || iteration.StageCount != 3 ||
		iteration.IRDigest == "" ||
		iteration.SourceDigest != history.LastCandidateSourceDigest ||
		iteration.FeedbackDigest != feedback.FeedbackDigest {
		t.Fatalf("unexpected self-improvement iteration: %#v", iteration)
	}
	if err := iteration.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestObserveRevisionSelfImprovementIterationRetainsSnapshotLinkFailure(t *testing.T) {
	source, history, feedback := selfImprovementIterationInputs(t)
	iteration, err := ObserveRevisionSelfImprovementIteration(source+"\n", history, feedback)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementIteration() error = nil, want source link failure")
	}
	if iteration.Status != "UNKNOWN" ||
		iteration.MissingStage != "revision-self-improvement-iteration-source-link" {
		t.Fatalf("unexpected unknown source-link iteration: %#v", iteration)
	}
}

func TestObserveRevisionSelfImprovementIterationRetainsHistoryFailure(t *testing.T) {
	source, history, feedback := selfImprovementIterationInputs(t)
	history.HistoryDigest = digestString("tampered")
	iteration, err := ObserveRevisionSelfImprovementIteration(source, history, feedback)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementIteration() error = nil, want history failure")
	}
	if iteration.Status != "UNKNOWN" ||
		iteration.MissingStage != "revision-self-improvement-iteration-history" {
		t.Fatalf("unexpected unknown history iteration: %#v", iteration)
	}
}

func TestObserveRevisionSelfImprovementIterationRetainsFeedbackLinkFailure(t *testing.T) {
	source, history, feedback := selfImprovementIterationInputs(t)
	feedback.WindowDigest = digestString("other-window")
	feedback.FeedbackDigest = digestRevisionSelfImprovementFeedback(feedback)
	iteration, err := ObserveRevisionSelfImprovementIteration(source, history, feedback)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementIteration() error = nil, want feedback link failure")
	}
	if iteration.Status != "UNKNOWN" ||
		iteration.MissingStage != "revision-self-improvement-iteration-feedback-link" {
		t.Fatalf("unexpected unknown feedback-link iteration: %#v", iteration)
	}
}

func TestObserveRevisionSelfImprovementIterationIsDeterministic(t *testing.T) {
	source, history, feedback := selfImprovementIterationInputs(t)
	first, err := ObserveRevisionSelfImprovementIteration(source, history, feedback)
	if err != nil {
		t.Fatalf("first ObserveRevisionSelfImprovementIteration() error = %v", err)
	}
	second, err := ObserveRevisionSelfImprovementIteration(source, history, feedback)
	if err != nil {
		t.Fatalf("second ObserveRevisionSelfImprovementIteration() error = %v", err)
	}
	if first.IterationDigest != second.IterationDigest {
		t.Fatal("same snapshot and history produced different iteration digest")
	}
}