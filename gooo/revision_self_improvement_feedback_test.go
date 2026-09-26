package gooo

import "testing"

func selfImprovementFeedbackWindow(t *testing.T, baselineBytes, candidateBytes int, irChanged bool) RevisionSelfImprovementWindow {
	t.Helper()
	baseline := selfImprovementReceiptWithProfile(t, baselineBytes, 1, false, true, true)
	candidate := selfImprovementReceiptWithProfile(t, candidateBytes, 1, irChanged, true, true)
	window, err := ObserveRevisionSelfImprovementWindow(baseline, candidate)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementWindow() error = %v", err)
	}
	return window
}

func TestObserveRevisionSelfImprovementFeedbackStable(t *testing.T) {
	window := selfImprovementFeedbackWindow(t, 4, 4, false)
	feedback, err := ObserveRevisionSelfImprovementFeedback(window)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementFeedback() error = %v", err)
	}
	if feedback.Status != "BOUND" || feedback.FeedbackSignal != "observe" ||
		feedback.FeedbackReason != "stable-observation" || !feedback.RequiresObservation ||
		feedback.RequiresReview || feedback.RequiresInspection || feedback.RequiresMeasurement {
		t.Fatalf("unexpected stable feedback: %#v", feedback)
	}
	if err := feedback.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestObserveRevisionSelfImprovementFeedbackNarrowerRequiresRemeasure(t *testing.T) {
	window := selfImprovementFeedbackWindow(t, 8, 2, false)
	feedback, err := ObserveRevisionSelfImprovementFeedback(window)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementFeedback() error = %v", err)
	}
	if feedback.ComparisonSignal != "narrower" || feedback.FeedbackSignal != "remeasure" ||
		!feedback.RequiresMeasurement || feedback.RequiresReview || feedback.RequiresInspection {
		t.Fatalf("unexpected narrower feedback: %#v", feedback)
	}
}

func TestObserveRevisionSelfImprovementFeedbackWiderRequiresReview(t *testing.T) {
	window := selfImprovementFeedbackWindow(t, 2, 8, false)
	feedback, err := ObserveRevisionSelfImprovementFeedback(window)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementFeedback() error = %v", err)
	}
	if feedback.ComparisonSignal != "wider" || feedback.FeedbackSignal != "review" ||
		!feedback.RequiresReview || feedback.RequiresMeasurement || feedback.RequiresInspection {
		t.Fatalf("unexpected wider feedback: %#v", feedback)
	}
}

func TestObserveRevisionSelfImprovementFeedbackMixedRequiresInspection(t *testing.T) {
	window := selfImprovementFeedbackWindow(t, 8, 2, true)
	feedback, err := ObserveRevisionSelfImprovementFeedback(window)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementFeedback() error = %v", err)
	}
	if feedback.ComparisonSignal != "mixed" || feedback.FeedbackSignal != "inspect" ||
		!feedback.RequiresInspection || !feedback.RequiresMeasurement ||
		feedback.RequiresReview {
		t.Fatalf("unexpected mixed feedback: %#v", feedback)
	}
}

func TestObserveRevisionSelfImprovementFeedbackRetainsWindowFailure(t *testing.T) {
	window := selfImprovementFeedbackWindow(t, 4, 4, false)
	window.WindowDigest = digestString("tampered")
	feedback, err := ObserveRevisionSelfImprovementFeedback(window)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementFeedback() error = nil, want window failure")
	}
	if feedback.Status != "UNKNOWN" || feedback.MissingStage != "revision-self-improvement-feedback-window" {
		t.Fatalf("unexpected unknown feedback: %#v", feedback)
	}
}

func TestObserveRevisionSelfImprovementFeedbackIsDeterministic(t *testing.T) {
	window := selfImprovementFeedbackWindow(t, 8, 2, false)
	first, err := ObserveRevisionSelfImprovementFeedback(window)
	if err != nil {
		t.Fatalf("first ObserveRevisionSelfImprovementFeedback() error = %v", err)
	}
	second, err := ObserveRevisionSelfImprovementFeedback(window)
	if err != nil {
		t.Fatalf("second ObserveRevisionSelfImprovementFeedback() error = %v", err)
	}
	if first.FeedbackDigest != second.FeedbackDigest {
		t.Fatal("same window produced different feedback digest")
	}
}