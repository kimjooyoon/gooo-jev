package gooo

import "testing"

func selfImprovementHistoryInputs(t *testing.T) ([]RevisionSelfImprovementWindow, []RevisionSelfImprovementFeedback) {
	t.Helper()
	windows := []RevisionSelfImprovementWindow{
		selfImprovementFeedbackWindow(t, 4, 4, false),
		selfImprovementFeedbackWindow(t, 8, 2, false),
		selfImprovementFeedbackWindow(t, 2, 8, false),
	}
	feedback := make([]RevisionSelfImprovementFeedback, 0, len(windows))
	for _, window := range windows {
		current, err := ObserveRevisionSelfImprovementFeedback(window)
		if err != nil {
			t.Fatalf("ObserveRevisionSelfImprovementFeedback() error = %v", err)
		}
		feedback = append(feedback, current)
	}
	return windows, feedback
}

func TestObserveRevisionSelfImprovementHistoryBindsOrderedSignals(t *testing.T) {
	windows, feedback := selfImprovementHistoryInputs(t)
	history, err := ObserveRevisionSelfImprovementHistory(windows, feedback)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementHistory() error = %v", err)
	}
	if history.Status != "BOUND" || history.ObservationCount != 3 ||
		history.HistorySignal != "mixed" {
		t.Fatalf("unexpected self-improvement history: %#v", history)
	}
	if history.FirstWindowDigest != windows[0].WindowDigest ||
		history.LastWindowDigest != windows[2].WindowDigest ||
		history.FirstFeedbackDigest != feedback[0].FeedbackDigest ||
		history.LastFeedbackDigest != feedback[2].FeedbackDigest {
		t.Fatalf("history boundaries were not retained: %#v", history)
	}
	if history.StableCount != 1 || history.NarrowerCount != 1 ||
		history.WiderCount != 1 || history.ObserveCount != 1 ||
		history.RemeasureCount != 1 || history.ReviewCount != 1 {
		t.Fatalf("history signal counts were not retained: %#v", history)
	}
	if err := history.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestObserveRevisionSelfImprovementHistoryRetainsFeedbackFailureIndex(t *testing.T) {
	windows, feedback := selfImprovementHistoryInputs(t)
	feedback[1].FeedbackDigest = digestString("tampered")
	history, err := ObserveRevisionSelfImprovementHistory(windows, feedback)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementHistory() error = nil, want feedback failure")
	}
	if history.Status != "UNKNOWN" ||
		history.MissingStage != "revision-self-improvement-history-feedback-1" {
		t.Fatalf("unexpected unknown history: %#v", history)
	}
}

func TestObserveRevisionSelfImprovementHistoryRetainsLengthMismatch(t *testing.T) {
	windows, feedback := selfImprovementHistoryInputs(t)
	feedback = feedback[:len(feedback)-1]
	history, err := ObserveRevisionSelfImprovementHistory(windows, feedback)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementHistory() error = nil, want count failure")
	}
	if history.Status != "UNKNOWN" ||
		history.MissingStage != "revision-self-improvement-history-feedback-count" {
		t.Fatalf("unexpected unknown count history: %#v", history)
	}
}

func TestObserveRevisionSelfImprovementHistoryIsDeterministic(t *testing.T) {
	windows, feedback := selfImprovementHistoryInputs(t)
	first, err := ObserveRevisionSelfImprovementHistory(windows, feedback)
	if err != nil {
		t.Fatalf("first ObserveRevisionSelfImprovementHistory() error = %v", err)
	}
	second, err := ObserveRevisionSelfImprovementHistory(windows, feedback)
	if err != nil {
		t.Fatalf("second ObserveRevisionSelfImprovementHistory() error = %v", err)
	}
	if first.HistoryDigest != second.HistoryDigest {
		t.Fatal("same ordered history produced different history digest")
	}
}