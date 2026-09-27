package gooo

import "testing"

func TestObserveCapabilityQueryFeedbackRequiresCompleteEvidenceForReview(t *testing.T) {
	response := DiscoverCapabilityQueryWithDeclaration(
		"What can gooo do with this declaration?",
		"entity Invoice\noperation reconcile\n",
	)
	feedback, err := ObserveCapabilityQueryFeedback(response, []string{"identity", "generation", "boundary", "receipt", "reverse_observation"})
	if err != nil {
		t.Fatalf("observe feedback: %v", err)
	}
	if err := feedback.Validate(); err != nil {
		t.Fatalf("validate feedback: %v", err)
	}
	if feedback.Disposition != CapabilityQueryFeedbackEligibleForReview {
		t.Fatalf("disposition = %q, want eligible review", feedback.Disposition)
	}
	if len(feedback.MissingStages) != 0 {
		t.Fatalf("missing stages = %#v", feedback.MissingStages)
	}
}

func TestObserveCapabilityQueryFeedbackPreservesUnknownAndRejectsUnknownStage(t *testing.T) {
	response := DiscoverCapabilityQuery("Can gooo infer an unsupported private fact?")
	feedback, err := ObserveCapabilityQueryFeedback(response, nil)
	if err != nil {
		t.Fatalf("observe unknown feedback: %v", err)
	}
	if err := feedback.Validate(); err != nil {
		t.Fatalf("validate unknown feedback: %v", err)
	}
	if feedback.Disposition != CapabilityQueryFeedbackPreserveUnknown {
		t.Fatalf("disposition = %q, want preserve unknown", feedback.Disposition)
	}
	if _, err := ObserveCapabilityQueryFeedback(response, []string{"invented_stage"}); err == nil {
		t.Fatal("unknown evidence stage was accepted")
	}
}
